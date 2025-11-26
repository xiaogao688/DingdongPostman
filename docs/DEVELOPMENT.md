# 开发指南

## 项目初始化

### 前置要求
- Go 1.21 或更高版本
- Make (可选，用于运行 Makefile 命令)

### 本地运行

#### 方式一：直接运行
```bash
go run main.go
```

#### 方式二：编译后运行
```bash
make build
./bin/dingdong-postman
```

#### 方式三：使用 Makefile
```bash
make run
```

### 环境变量

可以通过环境变量配置服务：

```bash
# 设置服务器端口 (默认: 8080)
export PORT=9090

# 设置运行环境 (默认: development)
export ENV=production

# 然后运行
go run main.go
```

### 可用的 Makefile 命令

```bash
# 编译项目
make build

# 运行项目
make run

# 运行测试
make test

# 清理编译产物
make clean

# 代码格式化
make fmt

# 代码检查 (需要安装 golangci-lint)
make lint

# Go vet 检查
make vet
```

### 项目结构

```
.
├── main.go                 # 应用入口
├── go.mod                  # Go 模块定义
├── Makefile               # 构建脚本
├── .gitignore             # Git 忽略规则
├── internal/
│   ├── config/           # 配置管理
│   ├── logger/           # 日志模块
│   └── server/           # HTTP 服务器
└── docs/
    └── DEVELOPMENT.md    # 本文件
```

### 测试服务

服务启动后，可以测试以下端点：

```bash
# 健康检查
curl http://localhost:8080/health

# 根路由
curl http://localhost:8080/
```

### 常见问题

**Q: 如何修改服务器端口？**
A: 设置 `PORT` 环境变量或修改 `internal/config/config.go` 中的默认值。

**Q: 如何添加新的路由？**
A: 在 `internal/server/server.go` 的 `setupRoutes()` 方法中添加新的路由处理器。

**Q: 如何添加依赖？**
A: 使用 `go get` 命令，例如：`go get github.com/some/package`

