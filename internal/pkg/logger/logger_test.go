package logger

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Level != "info" {
		t.Errorf("Expected level 'info', got '%s'", config.Level)
	}

	if !config.Console {
		t.Error("Expected console to be true")
	}

	if config.File == nil {
		t.Error("Expected file config to be not nil")
	}

	if !config.File.Compress {
		t.Error("Expected file compress to be true")
	}
}

// TestInitializerWithConsoleOnly 测试仅终端输出
func TestInitializerWithConsoleOnly(t *testing.T) {
	config := &Config{
		Level:   "info",
		Console: true,
		File: &FileConfig{
			Enabled: false,
		},
		Aliyun: &AliyunConfig{
			Enabled: false,
		},
	}

	initializer := NewInitializer(config)
	logger, err := initializer.Init()

	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Logger is nil")
	}

	// 测试日志方法
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")

	if err := logger.Sync(); err != nil {
		t.Fatalf("Failed to sync logger: %v", err)
	}
}

// TestInitializerWithFileOutput 测试文件输出
func TestInitializerWithFileOutput(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	logFile := tmpDir + "/test.log"

	config := &Config{
		Level:   "info",
		Console: false,
		File: &FileConfig{
			Enabled:    true,
			Path:       logFile,
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   false,
		},
		Aliyun: &AliyunConfig{
			Enabled: false,
		},
	}

	initializer := NewInitializer(config)
	logger, err := initializer.Init()

	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// 记录日志
	logger.Info("Test message", zap.String("key", "value"))

	if err := logger.Sync(); err != nil {
		t.Fatalf("Failed to sync logger: %v", err)
	}

	// 检查文件是否创建
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created at %s", logFile)
	}
}

// TestInitializerWithInvalidLevel 测试无效的日志级别
func TestInitializerWithInvalidLevel(t *testing.T) {
	config := &Config{
		Level:   "invalid",
		Console: true,
	}

	initializer := NewInitializer(config)
	_, err := initializer.Init()

	if err == nil {
		t.Error("Expected error for invalid log level")
	}
}

// TestConfigLoader 测试配置加载器
func TestConfigLoader(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := tmpDir + "/config.yaml"

	configContent := `logger:
  level: debug
  console: true
  file:
    enabled: true
    path: ./logs/app.log
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	loader := NewLoader(configFile)
	config, err := loader.Load()

	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.Level != "debug" {
		t.Errorf("Expected level 'debug', got '%s'", config.Level)
	}

	if !config.Console {
		t.Error("Expected console to be true")
	}

	if !config.File.Enabled {
		t.Error("Expected file to be enabled")
	}
}

// TestConfigLoaderWithMissingFile 测试配置文件不存在
func TestConfigLoaderWithMissingFile(t *testing.T) {
	loader := NewLoader("/nonexistent/config.yaml")
	config, err := loader.Load()

	// 应该返回默认配置，不报错
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if config == nil {
		t.Fatal("Config is nil")
	}

	// 应该是默认配置
	if config.Level != "info" {
		t.Errorf("Expected default level 'info', got '%s'", config.Level)
	}
}

// TestConfigLoaderWithEnvironmentVariables 测试环境变量覆盖
func TestConfigLoaderWithEnvironmentVariables(t *testing.T) {
	// 设置环境变量
	os.Setenv("ALIYUN_LOG_ENDPOINT", "test-endpoint")
	os.Setenv("ALIYUN_LOG_PROJECT", "test-project")
	os.Setenv("ALIYUN_LOG_REGION", "test-region")
	os.Setenv("ALIYUN_ACCESS_KEY_ID", "test-key-id")
	os.Setenv("ALIYUN_ACCESS_KEY_SECRET", "test-key-secret")

	defer func() {
		os.Unsetenv("ALIYUN_LOG_ENDPOINT")
		os.Unsetenv("ALIYUN_LOG_PROJECT")
		os.Unsetenv("ALIYUN_LOG_REGION")
		os.Unsetenv("ALIYUN_ACCESS_KEY_ID")
		os.Unsetenv("ALIYUN_ACCESS_KEY_SECRET")
	}()

	loader := NewLoader("/nonexistent/config.yaml")
	config, err := loader.Load()

	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.Aliyun.Endpoint != "test-endpoint" {
		t.Errorf("Expected endpoint 'test-endpoint', got '%s'", config.Aliyun.Endpoint)
	}

	if config.Aliyun.Project != "test-project" {
		t.Errorf("Expected project 'test-project', got '%s'", config.Aliyun.Project)
	}
}

// TestZapLoggerInterface 测试 Logger 接口实现
func TestZapLoggerInterface(t *testing.T) {
	config := &Config{
		Level:   "debug",
		Console: true,
	}

	initializer := NewInitializer(config)
	logger, err := initializer.Init()

	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// 验证 Logger 接口实现
	var _ Logger = logger

	// 测试所有方法
	logger.Debug("Debug", zap.String("test", "debug"))
	logger.Info("Info", zap.String("test", "info"))
	logger.Warn("Warn", zap.String("test", "warn"))
	logger.Error("Error", zap.String("test", "error"))

	if err := logger.Sync(); err != nil {
		t.Fatalf("Failed to sync: %v", err)
	}
}

// BenchmarkLogger 日志性能基准测试
func BenchmarkLogger(b *testing.B) {
	config := &Config{
		Level:   "info",
		Console: false,
		File: &FileConfig{
			Enabled: false,
		},
		Aliyun: &AliyunConfig{
			Enabled: false,
		},
	}

	initializer := NewInitializer(config)
	logger, _ := initializer.Init()
	defer logger.Sync()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message", zap.Int("iteration", i))
	}
}
