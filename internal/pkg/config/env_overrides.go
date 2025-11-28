package config

import (
	"os"
	"strconv"
	"strings"
)

// applyEnvOverrides 使用环境变量覆盖配置（优先级高于 YAML 文件）
func applyEnvOverrides(cfg *AppConfig) {
	if cfg == nil {
		return
	}

	// App 基础信息
	if v := getenvTrim("APP_NAME"); v != "" {
		cfg.App.Name = v
	}
	if v := getenvTrim("APP_ENV"); v != "" {
		cfg.App.Env = v
	}
	if v := getenvTrim("APP_VERSION"); v != "" {
		cfg.App.Version = v
	}

	// Logger 顶层
	if v := getenvTrim("LOGGER_LEVEL"); v != "" {
		cfg.Logger.Level = strings.ToLower(v)
	}
	if v := getenvTrim("LOGGER_CONSOLE"); v != "" {
		if b, ok := parseBool(v); ok {
			cfg.Logger.Console = b
		}
	}

	// Logger 文件输出
	if cfg.Logger.File == nil {
		cfg.Logger.File = &LoggerFileConfig{}
	}
	if v := getenvTrim("LOGGER_FILE_ENABLED"); v != "" {
		if b, ok := parseBool(v); ok {
			cfg.Logger.File.Enabled = b
		}
	}
	if v := getenvTrim("LOGGER_FILE_PATH"); v != "" {
		cfg.Logger.File.Path = v
	}
	if v := getenvTrim("LOGGER_FILE_MAX_SIZE"); v != "" {
		if n, ok := parseInt(v); ok {
			cfg.Logger.File.MaxSize = n
		}
	}
	if v := getenvTrim("LOGGER_FILE_MAX_BACKUPS"); v != "" {
		if n, ok := parseInt(v); ok {
			cfg.Logger.File.MaxBackups = n
		}
	}
	if v := getenvTrim("LOGGER_FILE_MAX_AGE"); v != "" {
		if n, ok := parseInt(v); ok {
			cfg.Logger.File.MaxAge = n
		}
	}
	if v := getenvTrim("LOGGER_FILE_COMPRESS"); v != "" {
		if b, ok := parseBool(v); ok {
			cfg.Logger.File.Compress = b
		}
	}

	// Logger 阿里云 SLS 输出
	if cfg.Logger.Aliyun == nil {
		cfg.Logger.Aliyun = &LoggerAliyunConfig{}
	}
	if v := getenvTrim("LOGGER_ALIYUN_ENABLED"); v != "" {
		if b, ok := parseBool(v); ok {
			cfg.Logger.Aliyun.Enabled = b
		}
	}
	// 兼容通用环境变量命名（来自 env.sh）
	if v := getenvTrim("ALIYUN_LOG_ENDPOINT"); v != "" {
		cfg.Logger.Aliyun.Endpoint = v
	}
	if v := getenvTrim("ALIYUN_LOG_PROJECT"); v != "" {
		cfg.Logger.Aliyun.Project = v
	}
	if v := getenvTrim("ALIYUN_LOG_REGION"); v != "" {
		cfg.Logger.Aliyun.Region = v
	}
	if v := getenvTrim("ALIYUN_ACCESS_KEY_ID"); v != "" {
		cfg.Logger.Aliyun.AccessKeyID = v
	}
	if v := getenvTrim("ALIYUN_ACCESS_KEY_SECRET"); v != "" {
		cfg.Logger.Aliyun.AccessKeySecret = v
	}
	// 其余阿里云字段
	if v := getenvTrim("LOGGER_ALIYUN_LOGSTORE"); v != "" {
		cfg.Logger.Aliyun.Logstore = v
	}
	if v := getenvTrim("LOGGER_ALIYUN_TOPIC"); v != "" {
		cfg.Logger.Aliyun.Topic = v
	}
	if v := getenvTrim("LOGGER_ALIYUN_SOURCE"); v != "" {
		cfg.Logger.Aliyun.Source = v
	}
	if v := getenvTrim("LOGGER_ALIYUN_BATCH_SIZE"); v != "" {
		if n, ok := parseInt(v); ok {
			cfg.Logger.Aliyun.BatchSize = n
		}
	}
	if v := getenvTrim("LOGGER_ALIYUN_FLUSH_INTERVAL"); v != "" {
		if n, ok := parseInt(v); ok {
			cfg.Logger.Aliyun.FlushInterval = n
		}
	}
}

func getenvTrim(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func parseBool(v string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "t", "yes", "y", "on":
		return true, true
	case "0", "false", "f", "no", "n", "off":
		return false, true
	default:
		return false, false
	}
}

func parseInt(v string) (int, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 0, false
	}
	return n, true
}
