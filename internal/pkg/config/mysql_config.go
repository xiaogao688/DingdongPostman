package config

import "os"

// MySQLConfig MySQL 配置结构
type MySQLConfig struct {
	// Enabled 是否启用 MySQL
	Enabled bool `yaml:"enabled" mapstructure:"enabled" default:"false"`

	// Host MySQL 服务器地址
	Host string `yaml:"host" mapstructure:"host" default:"localhost"`

	// Port MySQL 服务器端口
	Port int `yaml:"port" mapstructure:"port" default:"3306"`

	// Username MySQL 用户名
	Username string `yaml:"username" mapstructure:"username" default:"root"`

	// Password MySQL 密码（可从环境变量 MYSQL_PASSWORD 读取）
	Password string `yaml:"password" mapstructure:"password" default:""`

	// PasswordEnvVar 密码环境变量名称
	PasswordEnvVar string `yaml:"password_env_var" mapstructure:"password_env_var" default:"MYSQL_PASSWORD"`

	// Database 数据库名称
	Database string `yaml:"database" mapstructure:"database" default:""`

	// MaxOpenConns 最大打开连接数
	MaxOpenConns int `yaml:"max_open_conns" mapstructure:"max_open_conns" default:"25"`

	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int `yaml:"max_idle_conns" mapstructure:"max_idle_conns" default:"5"`

	// ConnMaxLifetime 连接最大生命周期（秒）
	ConnMaxLifetime int `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime" default:"300"`

	// ConnMaxIdleTime 连接最大空闲时间（秒）
	ConnMaxIdleTime int `yaml:"conn_max_idle_time" mapstructure:"conn_max_idle_time" default:"60"`

	// ParseTime 是否解析时间类型
	ParseTime bool `yaml:"parse_time" mapstructure:"parse_time" default:"true"`

	// Charset 字符集
	Charset string `yaml:"charset" mapstructure:"charset" default:"utf8mb4"`

	// Loc 时区
	Loc string `yaml:"loc" mapstructure:"loc" default:"Local"`

	// LogLevel 日志级别 (silent, error, warn, info)
	LogLevel string `yaml:"log_level" mapstructure:"log_level" default:"warn"`

	// SlowThreshold 慢查询阈值（毫秒）
	SlowThreshold int `yaml:"slow_threshold" mapstructure:"slow_threshold" default:"200"`
}

// DefaultMySQLConfig 返回默认 MySQL 配置
func DefaultMySQLConfig() *MySQLConfig {
	return &MySQLConfig{
		Enabled:         false,
		Host:            "localhost",
		Port:            3306,
		Username:        "root",
		Password:        "",
		PasswordEnvVar:  "MYSQL_PASSWORD",
		Database:        "",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 300,
		ConnMaxIdleTime: 60,
		ParseTime:       true,
		Charset:         "utf8mb4",
		Loc:             "Local",
		LogLevel:        "warn",
		SlowThreshold:   200,
	}
}

// GetPassword 获取 MySQL 密码，优先从环境变量读取
func (c *MySQLConfig) GetPassword() string {
	// 优先从环境变量读取
	if c.PasswordEnvVar != "" {
		if envPassword := os.Getenv(c.PasswordEnvVar); envPassword != "" {
			return envPassword
		}
	}
	// 其次使用配置文件中的密码
	return c.Password
}
