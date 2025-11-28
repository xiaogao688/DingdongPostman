package main

import (
	"fmt"
	"os"
	"time"

	appConfig "github.com/dingdong-postman/internal/pkg/config"
	appLogger "github.com/dingdong-postman/internal/pkg/logger"
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

	// 3) 示例输出
	log.Info("应用启动成功,shi", zap.Time("start_time", time.Now()))
	log.Error("这是一个错误示范",
		zap.String("str1", "str1----xxxx"),
		zap.String("str2", "str2----xxxx"),
		zap.Bool("bool", true))
	fmt.Printf("App: %s | Env: %s | Version: %s\n", cfg.App.Name, cfg.App.Env, cfg.App.Version)
}
