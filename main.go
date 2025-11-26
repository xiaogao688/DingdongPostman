package main

import (
	"github.com/dingdong-postman/internal/config"
	"github.com/dingdong-postman/internal/logger"
	"github.com/dingdong-postman/internal/server"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化日志
	log := logger.New()
	log.Info("叮咚邮差 (Dingdong Postman) - 分布式通知平台")
	log.Info("环境: " + cfg.Env)

	// 创建并启动服务器
	srv := server.New(cfg.Port, log)
	if err := srv.Start(); err != nil {
		log.Error("服务器启动失败: " + err.Error())
	}
}
