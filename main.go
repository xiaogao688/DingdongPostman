package main

import (
	"context"
	"fmt"
	"os"
	"time"

	appConfig "github.com/dingdong-postman/internal/pkg/config"
	appLogger "github.com/dingdong-postman/internal/pkg/logger"
	appMySQL "github.com/dingdong-postman/internal/pkg/mysql"
	appRedis "github.com/dingdong-postman/internal/pkg/redis"
	"go.uber.org/zap"
)

func main() {
	// 1) 初始化配置（支持通过环境变量 CONFIG_FILE 指定配置路径）
	configPath := os.Getenv("CONFIG_FILE")
	cfg, err := appConfig.Init(configPath)
	if err != nil {
		fmt.Printf("初始化配置失败: %v\n", err)
		os.Exit(1)
	}

	// 2) 初始化全局日志（基于配置）
	loggerCfg := appLogger.FromConfigLoggerConfig(&cfg.Logger)
	if err := appLogger.InitGlobal(loggerCfg); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = appLogger.GetGlobal().Sync()
	}()

	// 获取全局 Logger 实例
	log := appLogger.GetGlobal()

	// 3) 初始化全局 Redis 客户端
	if cfg.Redis.Enabled {
		if err := appRedis.InitGlobal(&cfg.Redis); err != nil {
			log.Warn("初始化 Redis 失败", zap.Error(err))
		} else {
			defer func() {
				_ = appRedis.Close()
			}()
			log.Info("Redis 初始化成功", zap.String("addr", cfg.Redis.Addr))
		}
	}

	// 4) 初始化全局 MySQL 数据库连接
	if cfg.MySQL.Enabled {
		if err := appMySQL.InitGlobal(&cfg.MySQL); err != nil {
			log.Warn("初始化 MySQL 失败", zap.Error(err))
		} else {
			defer func() {
				_ = appMySQL.Close()
			}()
			log.Info("MySQL 初始化成功", zap.String("host", cfg.MySQL.Host), zap.Int("port", cfg.MySQL.Port), zap.String("database", cfg.MySQL.Database))
		}
	}

	fmt.Printf("App: %s | Env: %s | Version: %s\n", cfg.App.Name, cfg.App.Env, cfg.App.Version)

	// 5) Redis 使用示例（如果已初始化）
	const timeout = 60 * time.Second
	if cfg.Redis.Enabled {
		err = appRedis.GetGlobal().Set(context.Background(), "test_key", "test_value", timeout)
		if err != nil {
			log.Error("Redis Set 失败", zap.Error(err))
		} else {
			log.Info("Redis Set 成功", zap.String("key", "test_key"), zap.String("value", "test_value"))
		}

		// 获取 Redis 值
		getValue, err := appRedis.GetGlobal().Get(context.Background(), "test_key")
		if err != nil {
			log.Error("Redis Get 失败", zap.Error(err))
		} else {
			log.Info("Redis Get 成功", zap.String("key", "test_key"), zap.String("value", getValue))
		}
	}

	// 6) MySQL 使用示例（如果已初始化）
	if cfg.MySQL.Enabled {
		db := appMySQL.GetGlobal()
		if db != nil {
			// 测试数据库连接
			if err := db.WithContext(context.Background()).Exec("SELECT 1").Error; err != nil {
				log.Error("MySQL 查询失败", zap.Error(err))
			} else {
				log.Info("MySQL 连接测试成功")
			}
		}
	}
}
