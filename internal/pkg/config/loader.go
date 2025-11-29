package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	cfgGlobal *AppConfig
	cfgOnce   sync.Once
	cfgErr    error
)

// Init 初始化全局配置。path 为空时，默认使用 "config/config.yaml"。
// 多次调用仅第一次生效，后续返回同一实例。
func Init(path string) (*AppConfig, error) {
	cfgOnce.Do(func() {
		if path == "" {
			path = "config/config.yaml"
		}
		var err error
		cfgGlobal, err = Load(path)
		if err != nil {
			cfgErr = err
			return
		}
		if err = cfgGlobal.Validate(); err != nil {
			cfgErr = err
			return
		}
	})
	return cfgGlobal, cfgErr
}

// Get 返回已初始化的全局配置，未初始化时返回 nil。
func Get() *AppConfig { return cfgGlobal }

// Load 使用 Viper 从指定路径加载配置，支持 YAML + 环境变量覆盖
func Load(path string) (*AppConfig, error) {
	def := Default()

	v := viper.New()
	v.SetConfigType("yaml")

	// 统一 key 风格：自动环境变量，分隔符替换
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// 设置默认值（来自 Default()）
	setDefaultsFrom(v, def)

	// 解析路径
	if path == "" {
		path = "config/config.yaml"
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("解析配置路径失败: %w", err)
	}

	// 配置文件：既支持显式文件，也支持目录+文件名
	if strings.HasSuffix(strings.ToLower(abs), ".yml") || strings.HasSuffix(strings.ToLower(abs), ".yaml") {
		v.SetConfigFile(abs)
	} else {
		v.AddConfigPath(abs)
		v.SetConfigName("config")
	}
	// 读取配置文件（不存在时忽略）
	_ = v.ReadInConfig()

	// 显式绑定环境变量到键（包含兼容你现有的命名）
	bindEnvKeys(v)

	// 反序列化
	cfg := &AppConfig{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("反序列化配置失败: %w", err)
	}

	return cfg, nil
}

func setDefaultsFrom(v *viper.Viper, def *AppConfig) {
	if def == nil {
		return
	}
	v.SetDefault("app.name", def.App.Name)
	v.SetDefault("app.env", def.App.Env)
	v.SetDefault("app.version", def.App.Version)

	v.SetDefault("logger.level", def.Logger.Level)
	v.SetDefault("logger.console", def.Logger.Console)

	if def.Logger.File != nil {
		v.SetDefault("logger.file.enabled", def.Logger.File.Enabled)
		v.SetDefault("logger.file.path", def.Logger.File.Path)
		v.SetDefault("logger.file.max_size", def.Logger.File.MaxSize)
		v.SetDefault("logger.file.max_backups", def.Logger.File.MaxBackups)
		v.SetDefault("logger.file.max_age", def.Logger.File.MaxAge)
		v.SetDefault("logger.file.compress", def.Logger.File.Compress)
	}
	if def.Logger.Aliyun != nil {
		v.SetDefault("logger.aliyun.enabled", def.Logger.Aliyun.Enabled)
		v.SetDefault("logger.aliyun.endpoint", def.Logger.Aliyun.Endpoint)
		v.SetDefault("logger.aliyun.project", def.Logger.Aliyun.Project)
		v.SetDefault("logger.aliyun.logstore", def.Logger.Aliyun.Logstore)
		v.SetDefault("logger.aliyun.region", def.Logger.Aliyun.Region)
		v.SetDefault("logger.aliyun.access_key_id", def.Logger.Aliyun.AccessKeyID)
		v.SetDefault("logger.aliyun.access_key_secret", def.Logger.Aliyun.AccessKeySecret)
		v.SetDefault("logger.aliyun.topic", def.Logger.Aliyun.Topic)
		v.SetDefault("logger.aliyun.source", def.Logger.Aliyun.Source)
		v.SetDefault("logger.aliyun.batch_size", def.Logger.Aliyun.BatchSize)
		v.SetDefault("logger.aliyun.flush_interval", def.Logger.Aliyun.FlushInterval)
	}
}

// 注意：config 模块现在不依赖 logger 模块
// 日志配置结构定义在 config 模块中（logger_config.go）
// logger 模块通过 FromConfigLoggerConfig() 函数来转换配置

func bindEnvKeys(v *viper.Viper) {
	pairs := map[string]string{
		// App
		"app.name":    "APP_NAME",
		"app.env":     "APP_ENV",
		"app.version": "APP_VERSION",

		// Logger 顶层
		"logger.level":   "LOGGER_LEVEL",
		"logger.console": "LOGGER_CONSOLE",

		// Logger 文件
		"logger.file.enabled":     "LOGGER_FILE_ENABLED",
		"logger.file.path":        "LOGGER_FILE_PATH",
		"logger.file.max_size":    "LOGGER_FILE_MAX_SIZE",
		"logger.file.max_backups": "LOGGER_FILE_MAX_BACKUPS",
		"logger.file.max_age":     "LOGGER_FILE_MAX_AGE",
		"logger.file.compress":    "LOGGER_FILE_COMPRESS",

		// Logger 阿里云：兼容 env.sh 与 LOGGER_*
		"logger.aliyun.enabled":           "LOGGER_ALIYUN_ENABLED",
		"logger.aliyun.endpoint":          "ALIYUN_LOG_ENDPOINT",
		"logger.aliyun.project":           "ALIYUN_LOG_PROJECT",
		"logger.aliyun.region":            "ALIYUN_LOG_REGION",
		"logger.aliyun.access_key_id":     "ALIYUN_ACCESS_KEY_ID",
		"logger.aliyun.access_key_secret": "ALIYUN_ACCESS_KEY_SECRET",
		"logger.aliyun.logstore":          "LOGGER_ALIYUN_LOGSTORE",
		"logger.aliyun.topic":             "LOGGER_ALIYUN_TOPIC",
		"logger.aliyun.source":            "LOGGER_ALIYUN_SOURCE",
		"logger.aliyun.batch_size":        "LOGGER_ALIYUN_BATCH_SIZE",
		"logger.aliyun.flush_interval":    "LOGGER_ALIYUN_FLUSH_INTERVAL",

		// Redis配置
		"redis.password": "REDIS_PASSWORD",
	}
	for key, env := range pairs {
		_ = v.BindEnv(key, env)
	}
}
