package mysql

import (
	"context"
	"fmt"
	"time"

	appLogger "github.com/dingdong-postman/internal/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// GormLogger 是 GORM 日志接口的实现，使用应用的日志模块
type GormLogger struct {
	logger appLogger.Logger
	level  logger.LogLevel
}

// NewGormLogger 创建一个新的 GORM 日志记录器
// 使用应用的全局日志模块来记录 GORM 的日志
func NewGormLogger(appLog appLogger.Logger, level logger.LogLevel) logger.Interface {
	if appLog == nil {
		return logger.Default
	}
	return &GormLogger{
		logger: appLog,
		level:  level,
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.level = level
	return l
}

// Info 记录信息级别日志
func (l *GormLogger) Info(_ context.Context, msg string, data ...interface{}) {
	if l.level < logger.Info {
		return
	}
	l.logger.Info(fmt.Sprintf(msg, data...))
}

// Warn 记录警告级别日志
func (l *GormLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.level < logger.Warn {
		return
	}
	l.logger.Warn(fmt.Sprintf(msg, data...))
}

// Error 记录错误级别日志
func (l *GormLogger) Error(_ context.Context, msg string, data ...interface{}) {
	if l.level < logger.Error {
		return
	}
	l.logger.Error(fmt.Sprintf(msg, data...))
}

// Trace 记录 SQL 执行日志
func (l *GormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.level >= logger.Error:
		// 错误日志
		l.logger.Error(
			"SQL 执行失败",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
			zap.Error(err),
		)

	case elapsed > time.Second && l.level >= logger.Warn:
		// 慢查询日志
		l.logger.Warn(
			"慢查询",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)

	case l.level >= logger.Info:
		// 普通日志
		l.logger.Info(
			"SQL 执行",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)
	}
}

// ParseGormLogLevel 将字符串转换为 GORM 日志级别
func ParseGormLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}
