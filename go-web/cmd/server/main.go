package main

import (
	"fmt"
	"go-web/pkg/logger"
	"go-web/pkg/setting"
	"go-web/internal/dao"
	"go-web/internal/routers"
	"github.com/spf13/viper"
)

func main() {
	// 先加载配置
	if err := setting.LoadConfig(); err != nil {
		fmt.Printf("配置加载失败: %v\n", err)
		return
	}

	// 初始化日志（使用配置）
	if err := logger.InitLogger(nil); err != nil {
		// 如果日志初始化失败，直接使用fmt输出到控制台
		fmt.Printf("日志初始化失败: %v\n", err)
		return
	}
	defer logger.Sync()

	logger.Info("开始初始化应用...")

	// 初始化数据库
	if err := dao.InitDB(); err != nil {
		logger.Fatalf("数据库初始化失败: %v", err)
	}
	
	// 程序退出时关闭数据库连接
	defer func() {
		if err := dao.CloseDB(); err != nil {
			logger.Errorw("关闭数据库连接失败", "error", err)
		} else {
			logger.Info("数据库连接已关闭")
		}
	}()

	logger.Info("数据库初始化成功")

	// 设置路由
	engine := routers.SetupRouter()
	logger.Info("路由设置完成")
	
	// 启动服务器
	host := viper.GetString("Server.Host")
	port := viper.GetInt("Server.Port")
	
	logger.Infow("服务器启动信息", 
		"host", host, 
		"port", port,
		"app_name", viper.GetString("Server"))
	
	if err := engine.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		logger.Fatalf("服务器启动失败: %v", err)
	}
}

