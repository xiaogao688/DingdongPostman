package redis

import (
	"sync"

	"github.com/dingdong-postman/internal/pkg/config"
)

var (
	globalClient Client
	clientOnce   sync.Once
)

// InitGlobal 初始化全局 Redis 客户端（仅初始化一次）
// 如需自定义配置，请在调用前构造好 *config.RedisConfig
func InitGlobal(cfg *config.RedisConfig) error {
	var initErr error
	clientOnce.Do(func() {
		initializer := NewInitializer(cfg)
		client, err := initializer.Init()
		if err != nil {
			initErr = err
			return
		}
		globalClient = client
	})
	return initErr
}

// GetGlobal 获取全局 Redis 客户端；若未初始化，返回 nil
func GetGlobal() Client {
	return globalClient
}

// MustGetGlobal 获取全局 Redis 客户端；未初始化则 panic（可用于必须已初始化的场景）
func MustGetGlobal() Client {
	if globalClient == nil {
		panic("redis: global client not initialized, call InitGlobal first")
	}
	return globalClient
}

// Close 关闭全局 Redis 客户端
func Close() error {
	if globalClient != nil {
		return globalClient.Close()
	}
	return nil
}
