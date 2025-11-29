package config

import "os"

// RedisConfig Redis 配置结构
type RedisConfig struct {
	// Enabled 是否启用 Redis
	Enabled bool `yaml:"enabled" mapstructure:"enabled" default:"false"`

	// Addr Redis 服务器地址 (host:port)
	Addr string `yaml:"addr" mapstructure:"addr" default:"localhost:6379"`

	// Username Redis 用户名（Redis 6.0+ 支持 ACL）
	Username string `yaml:"username" mapstructure:"username" default:""`

	// Password Redis 密码（可从环境变量 REDIS_PASSWORD 读取）
	Password string `yaml:"password" mapstructure:"password" default:""`

	// PasswordEnvVar 密码环境变量名称
	PasswordEnvVar string `yaml:"password_env_var" mapstructure:"password_env_var" default:"REDIS_PASSWORD"`

	// DB Redis 数据库编号
	DB int `yaml:"db" mapstructure:"db" default:"0"`

	// MaxRetries 最大重试次数
	MaxRetries int `yaml:"max_retries" mapstructure:"max_retries" default:"3"`

	// PoolSize 连接池大小
	PoolSize int `yaml:"pool_size" mapstructure:"pool_size" default:"10"`

	// DialTimeout 连接超时时间（秒）
	DialTimeout int `yaml:"dial_timeout" mapstructure:"dial_timeout" default:"5"`

	// ReadTimeout 读超时时间（秒）
	ReadTimeout int `yaml:"read_timeout" mapstructure:"read_timeout" default:"3"`

	// WriteTimeout 写超时时间（秒）
	WriteTimeout int `yaml:"write_timeout" mapstructure:"write_timeout" default:"3"`
}

// DefaultRedisConfig 返回默认 Redis 配置
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Enabled:        false,
		Addr:           "localhost:6379",
		Username:       "",
		Password:       "",
		PasswordEnvVar: "REDIS_PASSWORD",
		DB:             0,
		MaxRetries:     3,
		PoolSize:       10,
		DialTimeout:    5,
		ReadTimeout:    3,
		WriteTimeout:   3,
	}
}

// GetPassword 获取 Redis 密码，优先从环境变量读取
func (c *RedisConfig) GetPassword() string {
	// 优先从环境变量读取
	if c.PasswordEnvVar != "" {
		if envPassword := os.Getenv(c.PasswordEnvVar); envPassword != "" {
			return envPassword
		}
	}
	// 其次使用配置文件中的密码
	return c.Password
}
