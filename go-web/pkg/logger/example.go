package logger

import (
	"errors"
	"time"
	"go.uber.org/zap"
)

// ExampleUsage 展示日志封装的使用示例
func ExampleUsage() {
	// 1. 基本日志记录
	Info("应用启动成功")
	Debug("调试信息", zap.String("module", "example"))
	Warn("警告信息", zap.Int("count", 10))

	// 2. 格式化日志
	Infof("用户 %s 登录成功", "张三")
	Errorf("处理请求失败: %v", errors.New("连接超时"))

	// 3. 结构化日志（推荐）
	Infow("用户操作记录",
		"user_id", 123,
		"action", "login",
		"ip", "192.168.1.1",
		"timestamp", time.Now(),
	)

	// 4. 错误处理示例
	if err := someOperation(); err != nil {
		Errorw("操作失败",
			"error", err,
			"operation", "someOperation",
			"retry_count", 3,
		)
	}

	// 5. 性能监控示例
	start := time.Now()
	// 执行一些操作
	time.Sleep(100 * time.Millisecond)
	duration := time.Since(start)
	
	Infow("操作耗时统计",
		"operation", "data_processing",
		"duration_ms", duration.Milliseconds(),
		"status", "success",
	)
}

// 自定义配置示例
func ExampleWithCustomConfig() {
	config := &LogConfig{
		Level:      "debug",
		Format:     "console",
		Output:     "both",
		FilePath:   "logs/debug.log",
		MaxSize:    50,
		MaxBackups: 10,
		MaxAge:     3,
		Compress:   false,
	}

	if err := InitLogger(config); err != nil {
		panic(err)
	}
	defer Sync()

	// 使用自定义配置的日志记录
	Debug("调试模式已启用")
}

// 模拟操作函数
func someOperation() error {
	return errors.New("模拟错误")
}