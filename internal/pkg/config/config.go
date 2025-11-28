package config

import (
	"fmt"
	"os"
)

// AppSettings 应用基础信息
type AppSettings struct {
	Name    string `yaml:"name" mapstructure:"name"`
	Env     string `yaml:"env" mapstructure:"env"`         // development / staging / production
	Version string `yaml:"version" mapstructure:"version"` // 可选
}

// AppConfig 定义整个项目的配置根结构（对应 @config/config.yaml）
// 可按需在此文件中继续扩展你的业务配置结构体。
type AppConfig struct {
	// 应用基础信息
	App AppSettings `yaml:"app" mapstructure:"app"`

	// 日志配置（已实现：控制台、文件、阿里云日志服务）
	// 注意：日志配置定义在本模块中，避免循环依赖
	Logger LoggerConfig `yaml:"logger" mapstructure:"logger"`
}

// Default 返回项目的默认配置
func Default() *AppConfig {
	cfg := &AppConfig{}
	cfg.App.Name = "dingdong-postman"
	cfg.App.Env = envOrDefault("APP_ENV", "development")
	cfg.App.Version = "0.1.0"

	cfg.Logger = *DefaultLoggerConfig()
	return cfg
}

// envOrDefault 读取环境变量，不存在时返回默认值
func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

// Validate 对关键配置进行基础校验
func (c *AppConfig) Validate() error {
	// 校验日志配置
	// 若启用阿里云日志，校验必要字段（具体严格校验在 logger 初始化阶段也会再次进行）
	if c.Logger.Aliyun != nil && c.Logger.Aliyun.Enabled {
		if c.Logger.Aliyun.Logstore == "" {
			return fmt.Errorf("logger.aliyun.logstore 不能为空（已启用 aliyun）")
		}
	}
	return nil
}

// ToLoggerConfig 将 AppConfig 中的日志配置转换为 logger 模块的配置
// 这个方法用于在 logger 模块中使用配置
func (c *AppConfig) ToLoggerConfig() *LoggerConfig {
	return &c.Logger
}
