package mysql

import (
	"sync"

	"github.com/dingdong-postman/internal/pkg/config"
	appLogger "github.com/dingdong-postman/internal/pkg/logger"
	"gorm.io/gorm"
)

var (
	globalDB *gorm.DB
	dbOnce   sync.Once
)

// InitGlobal 初始化全局 MySQL 数据库连接（仅初始化一次）
// 如需自定义配置，请在调用前构造好 *config.MySQLConfig
// 日志记录器会自动使用应用的全局日志模块
func InitGlobal(cfg *config.MySQLConfig) error {
	var initErr error
	dbOnce.Do(func() {
		initializer := NewInitializer(cfg)
		db, err := initializer.Init()
		if err != nil {
			initErr = err
			return
		}
		globalDB = db
	})
	return initErr
}

// InitGlobalWithLogger 初始化全局 MySQL 数据库连接，并指定日志记录器
// 如需自定义日志记录器，请使用此函数
func InitGlobalWithLogger(cfg *config.MySQLConfig, logger appLogger.Logger) error {
	var initErr error
	dbOnce.Do(func() {
		initializer := NewInitializerWithLogger(cfg, logger)
		db, err := initializer.Init()
		if err != nil {
			initErr = err
			return
		}
		globalDB = db
	})
	return initErr
}

// GetGlobal 获取全局 MySQL 数据库连接；若未初始化，返回 nil
func GetGlobal() *gorm.DB {
	return globalDB
}

// MustGetGlobal 获取全局 MySQL 数据库连接；未初始化则 panic（可用于必须已初始化的场景）
func MustGetGlobal() *gorm.DB {
	if globalDB == nil {
		panic("mysql: global db not initialized, call InitGlobal first")
	}
	return globalDB
}

// Close 关闭全局 MySQL 数据库连接
func Close() error {
	if globalDB != nil {
		sqlDB, err := globalDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
