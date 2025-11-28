package logger

import (
	appConfig "github.com/dingdong-postman/internal/pkg/config"
)

// Config 日志配置结构
type Config struct {
	Level   string        `yaml:"level" mapstructure:"level" default:"info"`
	Console bool          `yaml:"console" mapstructure:"console" default:"true"`
	File    *FileConfig   `yaml:"file" mapstructure:"file"`
	Aliyun  *AliyunConfig `yaml:"aliyun" mapstructure:"aliyun"`
}

// FileConfig 文件日志配置
type FileConfig struct {
	Enabled    bool   `yaml:"enabled" mapstructure:"enabled" default:"false"`
	Path       string `yaml:"path" mapstructure:"path" default:"./logs/app.log"`
	MaxSize    int    `yaml:"max_size" mapstructure:"max_size" default:"100"`
	MaxBackups int    `yaml:"max_backups" mapstructure:"max_backups" default:"10"`
	MaxAge     int    `yaml:"max_age" mapstructure:"max_age" default:"30"`
	Compress   bool   `yaml:"compress" mapstructure:"compress" default:"true"`
}

// AliyunConfig 阿里云日志服务配置
type AliyunConfig struct {
	Enabled         bool   `yaml:"enabled" mapstructure:"enabled" default:"false"`
	Endpoint        string `yaml:"endpoint" mapstructure:"endpoint"`
	Project         string `yaml:"project" mapstructure:"project"`
	Logstore        string `yaml:"logstore" mapstructure:"logstore"`
	Region          string `yaml:"region" mapstructure:"region"`
	AccessKeyID     string `yaml:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret" mapstructure:"access_key_secret"`
	Topic           string `yaml:"topic" mapstructure:"topic" default:"app-log"`
	Source          string `yaml:"source" mapstructure:"source" default:"localhost"`
	BatchSize       int    `yaml:"batch_size" mapstructure:"batch_size" default:"100"`
	FlushInterval   int    `yaml:"flush_interval" mapstructure:"flush_interval" default:"5"`
}

const (
	defaultFileMaxSizeMB       = 100
	defaultFileMaxBackups      = 10
	defaultFileMaxAgeDays      = 30
	defaultAliyunBatchSize     = 100
	defaultAliyunFlushInterval = 5
)

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:   "info",
		Console: true,
		File: &FileConfig{
			Enabled:    false,
			Path:       "./logs/app.log",
			MaxSize:    defaultFileMaxSizeMB,
			MaxBackups: defaultFileMaxBackups,
			MaxAge:     defaultFileMaxAgeDays,
			Compress:   true,
		},
		Aliyun: &AliyunConfig{
			Enabled:       false,
			Topic:         "app-log",
			Source:        "localhost",
			BatchSize:     defaultAliyunBatchSize,
			FlushInterval: defaultAliyunFlushInterval,
		},
	}
}

// FromConfigLoggerConfig 将 config 模块的 LoggerConfig 转换为 logger 模块的 Config
func FromConfigLoggerConfig(c *appConfig.LoggerConfig) *Config {
	if c == nil {
		return DefaultConfig()
	}

	out := &Config{
		Level:   c.Level,
		Console: c.Console,
	}

	if c.File != nil {
		out.File = &FileConfig{
			Enabled:    c.File.Enabled,
			Path:       c.File.Path,
			MaxSize:    c.File.MaxSize,
			MaxBackups: c.File.MaxBackups,
			MaxAge:     c.File.MaxAge,
			Compress:   c.File.Compress,
		}
	}
	if c.Aliyun != nil {
		out.Aliyun = &AliyunConfig{
			Enabled:         c.Aliyun.Enabled,
			Endpoint:        c.Aliyun.Endpoint,
			Project:         c.Aliyun.Project,
			Logstore:        c.Aliyun.Logstore,
			Region:          c.Aliyun.Region,
			AccessKeyID:     c.Aliyun.AccessKeyID,
			AccessKeySecret: c.Aliyun.AccessKeySecret,
			Topic:           c.Aliyun.Topic,
			Source:          c.Aliyun.Source,
			BatchSize:       c.Aliyun.BatchSize,
			FlushInterval:   c.Aliyun.FlushInterval,
		}
	}
	return out
}
