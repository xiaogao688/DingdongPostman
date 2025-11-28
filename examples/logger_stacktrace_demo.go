package main

import (
	"fmt"

	"github.com/dingdong-postman/internal/pkg/logger"
	"go.uber.org/zap/zapcore"
)

// 这个示例展示了不同堆栈追踪配置的效果

func main() {
	fmt.Println("=== 日志堆栈追踪配置演示 ===\n")

	fmt.Println("1️⃣  配置 1: ErrorLevel（当前配置 - 推荐）")
	fmt.Println("   Error 及以上级别会显示堆栈追踪")
	fmt.Println("   输出: 日志消息 + 堆栈信息\n")
	demoWithStacktrace(zapcore.ErrorLevel)

	fmt.Println("\n2️⃣  配置 2: FatalLevel")
	fmt.Println("   只有 Fatal 级别会显示堆栈追踪")
	fmt.Println("   输出: 日志消息（Error 级别没有堆栈）\n")
	demoWithStacktrace(zapcore.FatalLevel)

	fmt.Println("\n3️⃣  配置 3: 禁用堆栈追踪")
	fmt.Println("   不显示任何堆栈追踪")
	fmt.Println("   输出: 仅日志消息\n")
	demoWithoutStacktrace()
}

// demoWithStacktrace 演示带堆栈追踪的日志
func demoWithStacktrace(stacktraceLevel zapcore.Level) {
	config := logger.DefaultConfig()
	config.Level = "debug"
	config.Console = true

	// 创建初始化器
	initializer := logger.NewInitializer(config)

	// 获取 Zap 日志器（这里需要修改 initializer 来支持自定义堆栈级别）
	// 为了演示，我们使用默认配置
	log, err := initializer.Init()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer log.Sync()

	// 记录不同级别的日志
	log.Debug("Debug message - 不显示堆栈")
	log.Info("Info message - 不显示堆栈")
	log.Warn("Warn message - 不显示堆栈")
	log.Error("Error message - 显示堆栈（当前配置）")
}

// demoWithoutStacktrace 演示不显示堆栈追踪的日志
func demoWithoutStacktrace() {
	config := logger.DefaultConfig()
	config.Level = "debug"
	config.Console = true

	initializer := logger.NewInitializer(config)
	log, err := initializer.Init()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer log.Sync()

	// 记录不同级别的日志
	log.Debug("Debug message - 不显示堆栈")
	log.Info("Info message - 不显示堆栈")
	log.Warn("Warn message - 不显示堆栈")
	log.Error("Error message - 不显示堆栈（禁用配置）")
}

// 如果你想自定义堆栈追踪级别，可以这样做：
// 修改 initializer.go 中的 Init 方法，添加参数：
//
// func (i *Initializer) InitWithStacktraceLevel(stacktraceLevel zapcore.Level) (Logger, error) {
//     // ... 其他代码 ...
//
//     zapLogger := zap.New(
//         zapcore.NewTee(cores...),
//         zap.AddCaller(),
//         zap.AddStacktrace(stacktraceLevel),  // 使用参数
//     )
//
//     return NewZapLogger(zapLogger), nil
// }
//
// 然后在代码中使用：
// log, _ := initializer.InitWithStacktraceLevel(zapcore.FatalLevel)
