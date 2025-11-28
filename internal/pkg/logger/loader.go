package logger

import (
	appConfig "github.com/dingdong-postman/internal/pkg/config"
)

// Loader 日志配置加载器（兼容保留）。
// 注意：为避免重复逻辑，内部委托给 config.Load，然后进行转换。
type Loader struct {
	configPath string
}

// NewLoader 创建一个新的配置加载器
func NewLoader(configPath string) *Loader {
	return &Loader{configPath: configPath}
}

// Load 加载日志配置：委托给 config.Load 再提取 logger 配置
func (l *Loader) Load() (*Config, error) {
	appCfg, err := appConfig.Load(l.configPath)
	if err != nil {
		return nil, err
	}
	return FromConfigLoggerConfig(&appCfg.Logger), nil
}
