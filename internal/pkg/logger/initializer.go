package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const defaultDirPerm os.FileMode = 0o755

// Initializer 日志初始化器，遵循 IoC 设计原则
type Initializer struct {
	config *Config
	mu     sync.Mutex
}

// NewInitializer 创建一个新的日志初始化器
func NewInitializer(config *Config) *Initializer {
	if config == nil {
		config = DefaultConfig()
	}
	return &Initializer{
		config: config,
	}
}

// Init 初始化日志系统
func (i *Initializer) Init() (Logger, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 获取日志级别
	level, err := zapcore.ParseLevel(i.config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心日志输出
	cores := []zapcore.Core{}

	// 1. 终端输出
	if i.config.Console {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 2. 文件输出
	if i.config.File != nil && i.config.File.Enabled {
		fileCore, err := i.createFileCore(encoderConfig, level)
		if err != nil {
			return nil, fmt.Errorf("failed to create file core: %w", err)
		}
		cores = append(cores, fileCore)
	}

	// 3. 阿里云日志服务输出
	if i.config.Aliyun != nil && i.config.Aliyun.Enabled {
		aliyunCore, err := i.createAliyunCore(encoderConfig, level)
		if err != nil {
			return nil, fmt.Errorf("failed to create aliyun core: %w", err)
		}
		cores = append(cores, aliyunCore)
	}

	// 如果没有配置任何输出，至少使用终端输出
	if len(cores) == 0 {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// 创建组合的日志记录器
	zapLogger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return NewZapLogger(zapLogger), nil
}

// createFileCore 创建文件日志核心
func (i *Initializer) createFileCore(encoderConfig zapcore.EncoderConfig, level zapcore.Level) (zapcore.Core, error) {
	// 创建日志目录
	dir := filepath.Dir(i.config.File.Path)
	if err := os.MkdirAll(dir, defaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// 使用 lumberjack 进行日志轮转
	writer := &lumberjack.Logger{
		Filename:   i.config.File.Path,
		MaxSize:    i.config.File.MaxSize,
		MaxBackups: i.config.File.MaxBackups,
		MaxAge:     i.config.File.MaxAge,
		Compress:   i.config.File.Compress,
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		level,
	), nil
}

// createAliyunCore 创建阿里云日志服务核心
func (i *Initializer) createAliyunCore(encoderConfig zapcore.EncoderConfig, level zapcore.Level) (zapcore.Core, error) {
	// 验证必要的配置
	if i.config.Aliyun.Endpoint == "" {
		return nil, fmt.Errorf("aliyun endpoint is required")
	}
	if i.config.Aliyun.Project == "" {
		return nil, fmt.Errorf("aliyun project is required")
	}
	if i.config.Aliyun.Logstore == "" {
		return nil, fmt.Errorf("aliyun logstore is required")
	}
	if i.config.Aliyun.AccessKeyID == "" {
		return nil, fmt.Errorf("aliyun access_key_id is required")
	}
	if i.config.Aliyun.AccessKeySecret == "" {
		return nil, fmt.Errorf("aliyun access_key_secret is required")
	}

	// 创建阿里云客户端
	provider := sls.NewStaticCredentialsProvider(
		i.config.Aliyun.AccessKeyID,
		i.config.Aliyun.AccessKeySecret,
		"",
	)

	client := sls.CreateNormalInterfaceV2(i.config.Aliyun.Endpoint, provider)
	client.SetAuthVersion(sls.AuthV4)
	if i.config.Aliyun.Region != "" {
		client.SetRegion(i.config.Aliyun.Region)
	}

	// 验证连接
	if _, err := client.ListProject(); err != nil {
		return nil, fmt.Errorf("failed to connect to aliyun log service: %w", err)
	}

	// 验证 Logstore 存在
	if _, err := client.GetLogStore(i.config.Aliyun.Project, i.config.Aliyun.Logstore); err != nil {
		return nil, fmt.Errorf("logstore %s does not exist: %w", i.config.Aliyun.Logstore, err)
	}

	// 创建阿里云日志写入器
	writer := NewAliyunWriter(
		client,
		i.config.Aliyun.Project,
		i.config.Aliyun.Logstore,
		i.config.Aliyun.Topic,
		i.config.Aliyun.Source,
		i.config.Aliyun.BatchSize,
		time.Duration(i.config.Aliyun.FlushInterval)*time.Second,
	)

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		level,
	), nil
}
