package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dingdong-postman/internal/logger"
)

type Server struct {
	port   int
	logger *logger.Logger
	mux    *http.ServeMux
}

func New(port int, logger *logger.Logger) *Server {
	return &Server{
		port:   port,
		logger: logger,
		mux:    http.NewServeMux(),
	}
}

func (s *Server) setupRoutes() {
	// 健康检查端点
	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// 根路由
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"叮咚邮差 - 分布式通知平台"}`))
	})
}

func (s *Server) Start() error {
	s.setupRoutes()

	addr := ":" + strconv.Itoa(s.port)
	s.logger.Info(fmt.Sprintf("服务器启动在 http://localhost:%d", s.port))

	return http.ListenAndServe(addr, s.mux)
}
