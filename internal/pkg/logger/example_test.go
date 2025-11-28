package logger

import (
	"fmt"
	"os"
)

// ExampleUsage 展示如何使用日志模块
func ExampleUsage() {
	// 方式 1: 使用配置加载器从 YAML 文件加载配置
	loader := NewLoader("config/config.yaml")
	config, err := loader.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// 方式 2: 手动创建配置
	// config := &Config{
	//     Level:   "info",
	//     Console: true,
	//     File: &FileConfig{
	//         Enabled: true,
	//         Path:    "./logs/app.log",
	//     },
	//     Aliyun: &AliyunConfig{
	//         Enabled:         true,
	//         Endpoint:        os.Getenv("ALIYUN_LOG_ENDPOINT"),
	//         Project:         os.Getenv("ALIYUN_LOG_PROJECT"),
	//         Logstore:        "my-logstore",
	//         Region:          os.Getenv("ALIYUN_LOG_REGION"),
	//         AccessKeyID:     os.Getenv("ALIYUN_ACCESS_KEY_ID"),
	//         AccessKeySecret: os.Getenv("ALIYUN_ACCESS_KEY_SECRET"),
	//     },
	// }

	// 创建初始化器
	initializer := NewInitializer(config)

	// 初始化日志系统
	logger, err := initializer.Init()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer logger.Sync()

	// 使用日志记录器
	logger.Info("Application started")
	logger.Debug("Debug message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// 输出示例
	// time=2024-11-26T17:36:27.093Z level=info msg="Application started"
	// time=2024-11-26T17:36:27.093Z level=debug msg="Debug message"
	// time=2024-11-26T17:36:27.093Z level=warn msg="Warning message"
	// time=2024-11-26T17:36:27.093Z level=error msg="Error message"
}

// ExampleWithEnvironmentVariables 展示如何使用环境变量配置
func ExampleWithEnvironmentVariables() {
	// 加载环境变量
	os.Setenv("ALIYUN_LOG_ENDPOINT", "cn-beijing.log.aliyuncs.com")
	os.Setenv("ALIYUN_LOG_PROJECT", "my-project")
	os.Setenv("ALIYUN_LOG_REGION", "cn-beijing")
	os.Setenv("ALIYUN_ACCESS_KEY_ID", "your-access-key-id")
	os.Setenv("ALIYUN_ACCESS_KEY_SECRET", "your-access-key-secret")

	// 加载配置
	loader := NewLoader("config/config.yaml")
	config, err := loader.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// 初始化日志系统
	initializer := NewInitializer(config)
	logger, err := initializer.Init()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer logger.Sync()

	// 使用日志记录器
	logger.Info("Logger with environment variables configured")
}
