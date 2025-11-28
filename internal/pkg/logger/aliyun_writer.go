package logger

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
)

// unixToUint32 returns a clamped uint32 representation of t.Unix().
// It prevents potential int64->uint32 overflows (gosec G115) by bounding the value.
func unixToUint32(t time.Time) uint32 {
	sec := t.Unix()
	if sec < 0 {
		return 0
	}
	// max uint32 value as int64
	const maxUint32 = ^uint32(0)
	if sec > int64(maxUint32) {
		return maxUint32
	}
	return uint32(sec)
}

// AliyunWriter 阿里云日志写入器
type AliyunWriter struct {
	client        sls.ClientInterface
	project       string
	logstore      string
	topic         string
	source        string
	batchSize     int
	flushInterval time.Duration

	mu       sync.Mutex
	buffer   []*sls.LogContent
	ticker   *time.Ticker
	stopChan chan struct{}
	done     chan struct{}
}

// NewAliyunWriter 创建一个新的阿里云日志写入器
func NewAliyunWriter(
	client sls.ClientInterface,
	project string,
	logstore string,
	topic string,
	source string,
	batchSize int,
	flushInterval time.Duration,
) *AliyunWriter {
	w := &AliyunWriter{
		client:        client,
		project:       project,
		logstore:      logstore,
		topic:         topic,
		source:        source,
		batchSize:     batchSize,
		flushInterval: flushInterval,
		buffer:        make([]*sls.LogContent, 0, batchSize),
		stopChan:      make(chan struct{}),
		done:          make(chan struct{}),
	}

	// 启动后台刷新 goroutine
	go w.flushLoop()

	return w
}

// Write 实现 io.Writer 接口
func (w *AliyunWriter) Write(p []byte) (n int, err error) {
	// 解析日志内容
	var logData map[string]interface{}
	if err := json.Unmarshal(p, &logData); err != nil {
		// 如果无法解析为 JSON，直接作为消息写入
		logData = map[string]interface{}{
			"msg": string(p),
		}
	}

	// 转换为阿里云日志格式
	logContents := make([]*sls.LogContent, 0, len(logData))
	for key, value := range logData {
		strValue := fmt.Sprintf("%v", value)
		logContents = append(logContents, &sls.LogContent{
			Key:   proto.String(key),
			Value: proto.String(strValue),
		})
	}

	w.mu.Lock()
	w.buffer = append(w.buffer, logContents...)

	// 如果缓冲区达到批量大小，立即刷新
	if len(w.buffer) >= w.batchSize {
		w.flushLocked()
	}
	w.mu.Unlock()

	return len(p), nil
}

// flushLoop 后台定时刷新日志
func (w *AliyunWriter) flushLoop() {
	w.ticker = time.NewTicker(w.flushInterval)
	defer w.ticker.Stop()

	for {
		select {
		case <-w.ticker.C:
			w.mu.Lock()
			if len(w.buffer) > 0 {
				w.flushLocked()
			}
			w.mu.Unlock()
		case <-w.stopChan:
			// 最后一次刷新
			w.mu.Lock()
			if len(w.buffer) > 0 {
				w.flushLocked()
			}
			w.mu.Unlock()
			close(w.done)
			return
		}
	}
}

// flushLocked 刷新缓冲区中的日志（需要在持有锁的情况下调用）
func (w *AliyunWriter) flushLocked() {
	if len(w.buffer) == 0 {
		return
	}

	// 创建日志组
	logs := make([]*sls.Log, 0, len(w.buffer)/10+1)
	currentLog := &sls.Log{
		Time:     proto.Uint32(unixToUint32(time.Now())),
		Contents: make([]*sls.LogContent, 0),
	}

	for i, content := range w.buffer {
		currentLog.Contents = append(currentLog.Contents, content)

		// 每 10 条日志创建一个新的 Log 对象（阿里云限制）
		if (i+1)%10 == 0 || i == len(w.buffer)-1 {
			logs = append(logs, currentLog)
			if i < len(w.buffer)-1 {
				currentLog = &sls.Log{
					Time:     proto.Uint32(unixToUint32(time.Now())),
					Contents: make([]*sls.LogContent, 0),
				}
			}
		}
	}

	logGroup := &sls.LogGroup{
		Topic:  proto.String(w.topic),
		Source: proto.String(w.source),
		Logs:   logs,
	}

	// 写入日志到阿里云
	if err := w.client.PutLogs(w.project, w.logstore, logGroup); err != nil {
		// 这里可以添加错误处理逻辑，例如记录到本地或重试
		fmt.Printf("failed to write logs to aliyun: %v\n", err)
	}

	// 清空缓冲区
	w.buffer = w.buffer[:0]
}

// Flush 手动刷新缓冲区
func (w *AliyunWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.flushLocked()
	return nil
}

// Sync 实现 zapcore.WriteSyncer 接口，使得 logger.Sync() 会触发一次 Flush
func (w *AliyunWriter) Sync() error {
	return w.Flush()
}

// Close 关闭写入器（用于优雅退出），会先做最后一次刷新
func (w *AliyunWriter) Close() error {
	close(w.stopChan)
	<-w.done
	return nil
}
