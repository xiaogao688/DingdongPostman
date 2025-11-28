package config

// LoggerConfig 日志配置结构（独立于 logger 模块）
// 这个结构定义在 config 模块中，避免 config 依赖 logger 模块
type LoggerConfig struct {
	// Level 日志级别 (debug, info, warn, error, fatal)
	Level string `yaml:"level" mapstructure:"level" default:"info"`

	// Console 是否输出到终端
	Console bool `yaml:"console" mapstructure:"console" default:"true"`

	// File 文件日志配置
	File *LoggerFileConfig `yaml:"file" mapstructure:"file"`

	// Aliyun 阿里云日志服务配置
	Aliyun *LoggerAliyunConfig `yaml:"aliyun" mapstructure:"aliyun"`
}

// LoggerFileConfig 文件日志配置
type LoggerFileConfig struct {
	// Enabled 是否启用文件日志
	Enabled bool `yaml:"enabled" mapstructure:"enabled" default:"false"`

	// Path 日志文件路径
	Path string `yaml:"path" mapstructure:"path" default:"./logs/app.log"`

	// MaxSize 单个日志文件最大大小（MB）
	MaxSize int `yaml:"max_size" mapstructure:"max_size" default:"100"`

	// MaxBackups 保留的最大备份文件数
	MaxBackups int `yaml:"max_backups" mapstructure:"max_backups" default:"10"`

	// MaxAge 日志文件最大保留天数
	MaxAge int `yaml:"max_age" mapstructure:"max_age" default:"30"`

	// Compress 是否压缩备份文件
	Compress bool `yaml:"compress" mapstructure:"compress" default:"true"`
}

// LoggerAliyunConfig 阿里云日志服务配置
type LoggerAliyunConfig struct {
	// Enabled 是否启用阿里云日志服务
	Enabled bool `yaml:"enabled" mapstructure:"enabled" default:"false"`

	// Endpoint 阿里云日志服务端点
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`

	// Project 项目名称
	Project string `yaml:"project" mapstructure:"project"`

	// Logstore Logstore 名称
	Logstore string `yaml:"logstore" mapstructure:"logstore"`

	// Region 地域信息
	Region string `yaml:"region" mapstructure:"region"`

	// AccessKeyID 访问密钥 ID
	AccessKeyID string `yaml:"access_key_id" mapstructure:"access_key_id"`

	// AccessKeySecret 访问密钥密码
	AccessKeySecret string `yaml:"access_key_secret" mapstructure:"access_key_secret"`

	// Topic 日志主题
	Topic string `yaml:"topic" mapstructure:"topic" default:"app-log"`

	// Source 日志来源
	Source string `yaml:"source" mapstructure:"source" default:"localhost"`

	// BatchSize 批量写入日志的大小
	BatchSize int `yaml:"batch_size" mapstructure:"batch_size" default:"100"`

	// FlushInterval 刷新间隔（秒）
	FlushInterval int `yaml:"flush_interval" mapstructure:"flush_interval" default:"5"`
}

// DefaultLoggerConfig 返回默认日志配置
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:   "info",
		Console: true,
		File: &LoggerFileConfig{
			Enabled:    false,
			Path:       "./logs/app.log",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
		},
		Aliyun: &LoggerAliyunConfig{
			Enabled:       false,
			Topic:         "app-log",
			Source:        "localhost",
			BatchSize:     100,
			FlushInterval: 5,
		},
	}
}
