package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	globalLogger Logger
	loggerOnce   sync.Once
)

// InitGlobal 初始化全局 Logger 实例（仅初始化一次）
// 如需自定义配置，请在调用前构造好 *Config
func InitGlobal(cfg *Config) error {
	var initErr error
	loggerOnce.Do(func() {
		initializer := NewInitializer(cfg)
		l, err := initializer.Init()
		if err != nil {
			initErr = err
			return
		}
		globalLogger = l
	})
	return initErr
}

// GetGlobal 获取全局 Logger 实例；若未初始化，返回一个 no-op Logger
func GetGlobal() Logger {
	if globalLogger == nil {
		return NewZapLogger(zap.NewNop())
	}
	return globalLogger
}

// MustGetGlobal 获取全局 Logger；未初始化则 panic（可用于必须已初始化的场景）
func MustGetGlobal() Logger {
	if globalLogger == nil {
		panic("logger: global logger not initialized, call InitGlobal first")
	}
	return globalLogger
}
