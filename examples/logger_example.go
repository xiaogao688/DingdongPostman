package main

import (
	"fmt"
	"os"

	"github.com/dingdong-postman/internal/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 示例 1: 基础使用 - 仅终端输出
	//fmt.Println("=== 示例 1: 基础使用 - 仅终端输出 ===")
	//basicExample()

	fmt.Println("\n=== 示例 2: 从配置文件加载 ===")
	configFileExample()
	//
	//fmt.Println("\n=== 示例 3: 手动配置所有输出 ===")
	//fullConfigExample()
}

// basicExample 基础使用示例
func basicExample() {
	// 创建默认配置
	config := logger.DefaultConfig()
	config.Level = "debug"
	config.Console = true

	// 初始化日志系统
	initializer := logger.NewInitializer(config)
	log, err := initializer.Init()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer log.Sync()

	// 使用日志
	log.Debug("This is a debug message")
	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")
}

// configFileExample 从配置文件加载示例
func configFileExample() {
	// 加载配置文件
	loader := logger.NewLoader("/home/gao/code/DingdongPostman/config/config.yaml")
	config, err := loader.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		// 使用默认配置
		config = logger.DefaultConfig()
	}

	// 初始化日志系统
	initializer := logger.NewInitializer(config)
	log, err := initializer.Init()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer log.Sync()

	// 使用日志
	log.Info("Logger initialized from config file")
	log.Info("Logging to console and file")
}

// fullConfigExample 完整配置示例
func fullConfigExample() {
	// 设置环境变量（演示用）
	os.Setenv("ALIYUN_LOG_ENDPOINT", "cn-beijing.log.aliyuncs.com")
	os.Setenv("ALIYUN_LOG_PROJECT", "demo-project")
	os.Setenv("ALIYUN_LOG_REGION", "cn-beijing")
	os.Setenv("ALIYUN_ACCESS_KEY_ID", "demo-key-id")
	os.Setenv("ALIYUN_ACCESS_KEY_SECRET", "demo-key-secret")

	// 创建完整配置
	config := &logger.Config{
		Level:   "info",
		Console: true,
		File: &logger.FileConfig{
			Enabled:    true,
			Path:       "./logs/app.log",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
		},
		// 注意: 阿里云配置需要真实的凭证才能正常工作
		// Aliyun: &logger.AliyunConfig{
		// 	Enabled:         true,
		// 	Endpoint:        os.Getenv("ALIYUN_LOG_ENDPOINT"),
		// 	Project:         os.Getenv("ALIYUN_LOG_PROJECT"),
		// 	Logstore:        "my-logstore",
		// 	Region:          os.Getenv("ALIYUN_LOG_REGION"),
		// 	AccessKeyID:     os.Getenv("ALIYUN_ACCESS_KEY_ID"),
		// 	AccessKeySecret: os.Getenv("ALIYUN_ACCESS_KEY_SECRET"),
		// 	Topic:           "app-log",
		// 	Source:          "localhost",
		// 	BatchSize:       100,
		// 	FlushInterval:   5,
		// },
	}

	// 初始化日志系统
	initializer := logger.NewInitializer(config)
	log, err := initializer.Init()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer log.Sync()

	// 使用日志
	log.Info("Application started",
		zap.String("version", "1.0.0"),
		zap.String("environment", "development"),
	)

	log.Debug("Debug information",
		zap.Int("request_id", 12345),
		zap.String("user_id", "user_123"),
	)

	log.Warn("Warning: High memory usage",
		zap.Float64("memory_percent", 85.5),
	)

	log.Error("Error processing request",
		zap.String("error", "connection timeout"),
		zap.Int("retry_count", 3),
	)

	fmt.Println("Logs have been written to console and file (./logs/app.log)")
}
