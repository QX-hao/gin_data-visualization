package main

import (
	"fmt"
	"log"
	"go-web/pkg/setting"
	"go-web/internal/dao"
	"go-web/internal/routers"
	"github.com/spf13/viper"
)

func main() {
	// 加载配置文件
	if err := setting.LoadConfig(); err != nil {
		log.Fatalf("配置文件加载失败: %v", err)
	}

	// 初始化数据库
	if err := dao.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	
	// 程序退出时关闭数据库连接
	defer func() {
		if err := dao.CloseDB(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		} else {
			log.Println("数据库连接已关闭")
		}
	}()

	// 设置路由
	engine := routers.SetupRouter()
	
	// 启动服务器
	log.Printf("服务器启动在: %s:%d", viper.GetString("AppName.Host"), viper.GetInt("AppName.Port"))
	if err := engine.Run(fmt.Sprintf("%s:%d", viper.GetString("AppName.Host"), viper.GetInt("AppName.Port"))); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

