package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/dingdong-postman/internal/pkg/config"
	appLogger "github.com/dingdong-postman/internal/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const TestMySQLConnectionTimeout = 5 * time.Second

// Initializer MySQL 初始化器
type Initializer struct {
	cfg    *config.MySQLConfig
	logger appLogger.Logger
}

// NewInitializer 创建 MySQL 初始化器
func NewInitializer(cfg *config.MySQLConfig) *Initializer {
	return &Initializer{
		cfg:    cfg,
		logger: appLogger.GetGlobal(),
	}
}

// NewInitializerWithLogger 创建 MySQL 初始化器，并指定日志记录器
func NewInitializerWithLogger(cfg *config.MySQLConfig, logger appLogger.Logger) *Initializer {
	if logger == nil {
		logger = appLogger.GetGlobal()
	}
	return &Initializer{
		cfg:    cfg,
		logger: logger,
	}
}

// Init 初始化 MySQL 数据库连接
func (i *Initializer) Init() (*gorm.DB, error) {
	if i.cfg == nil {
		return nil, fmt.Errorf("mysql config is nil")
	}

	if !i.cfg.Enabled {
		return nil, fmt.Errorf("mysql is not enabled")
	}

	// 获取密码（优先从环境变量读取）
	password := i.cfg.GetPassword()

	// 构建 DSN (Data Source Name)
	// 格式: username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		i.cfg.Username,
		password,
		i.cfg.Host,
		i.cfg.Port,
		i.cfg.Database,
		i.cfg.Charset,
		i.cfg.ParseTime,
		i.cfg.Loc,
	)

	// 转换日志级别
	logLevel := ParseGormLogLevel(i.cfg.LogLevel)

	// 创建 GORM 日志记录器（使用应用的日志模块）
	gormLogger := NewGormLogger(i.logger, logLevel)

	// 创建 GORM 数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	// 获取底层 SQL 数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql db: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxOpenConns(i.cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(i.cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(i.cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(i.cfg.ConnMaxIdleTime) * time.Second)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), TestMySQLConnectionTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return db, nil
}
