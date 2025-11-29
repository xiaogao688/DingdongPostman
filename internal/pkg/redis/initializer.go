package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/dingdong-postman/internal/pkg/config"
	"github.com/redis/go-redis/v9"
)

const TestRedisConnectionTimeout = 5 * time.Second

// Initializer Redis 初始化器
type Initializer struct {
	cfg *config.RedisConfig
}

// NewInitializer 创建 Redis 初始化器
func NewInitializer(cfg *config.RedisConfig) *Initializer {
	return &Initializer{
		cfg: cfg,
	}
}

// Init 初始化 Redis 客户端
func (i *Initializer) Init() (Client, error) {
	if i.cfg == nil {
		return nil, fmt.Errorf("redis config is nil")
	}

	if !i.cfg.Enabled {
		return nil, fmt.Errorf("redis is not enabled")
	}

	// 获取密码（优先从环境变量读取）
	password := i.cfg.GetPassword()

	// 创建 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:         i.cfg.Addr,
		Username:     i.cfg.Username,
		Password:     password,
		DB:           i.cfg.DB,
		MaxRetries:   i.cfg.MaxRetries,
		PoolSize:     i.cfg.PoolSize,
		DialTimeout:  time.Duration(i.cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(i.cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(i.cfg.WriteTimeout) * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), TestRedisConnectionTimeout)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return NewClient(redisClient), nil
}
