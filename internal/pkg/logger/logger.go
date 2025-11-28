package logger

import (
	"go.uber.org/zap"
)

// Logger 定义日志接口，遵循 IoC 设计原则
type Logger interface {
	// Debug 记录调试级别日志
	Debug(msg string, fields ...zap.Field)
	// Info 记录信息级别日志
	Info(msg string, fields ...zap.Field)
	// Warn 记录警告级别日志
	Warn(msg string, fields ...zap.Field)
	// Error 记录错误级别日志
	Error(msg string, fields ...zap.Field)
	// Fatal 记录致命错误日志并退出程序
	Fatal(msg string, fields ...zap.Field)
	// Sync 刷新日志缓冲区
	Sync() error
}

// zapLogger 是 Logger 接口的实现
type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger 创建一个新的 Zap 日志记录器
func NewZapLogger(logger *zap.Logger) Logger {
	return &zapLogger{
		logger: logger,
	}
}

func (z *zapLogger) Debug(msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}

func (z *zapLogger) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

func (z *zapLogger) Warn(msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}

func (z *zapLogger) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

func (z *zapLogger) Fatal(msg string, fields ...zap.Field) {
	z.logger.Fatal(msg, fields...)
}

func (z *zapLogger) Sync() error {
	return z.logger.Sync()
}
