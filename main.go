package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"gin_data-visualization/pkg/config"
	"gin_data-visualization/router"
	"gin_data-visualization/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 应用初始化
	cfg, db := initializeApp()
	
	// 设置优雅关闭
	setupGracefulShutdown(db)
	
	// 启动服务器
	if err := startServer(cfg, db); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// setupGracefulShutdown 设置优雅关闭
func setupGracefulShutdown(db *gorm.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	// 异步等待关闭信号
	go func() {
		<-quit
		log.Println("接收到关闭信号，正在优雅关闭...")
		
		closeDatabase(db)
		log.Println("服务器已关闭")
		os.Exit(0)
	}()
}

// initializeApp 应用初始化
func initializeApp() (*config.Config, *gorm.DB) {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Printf("配置文件加载失败，使用默认配置: %v", err)
	}
	
	cfg := config.GetConfig()
	logConfig(cfg)
	
	// 检查系统要求
	if err := utils.CheckSystemRequirements(); err != nil {
		log.Fatalf("系统要求检查失败: %v", err)
	}
	
	// 初始化组件
	db, err := utils.InitAll()
	if err != nil {
		log.Fatalf("组件初始化失败: %v", err)
	}
	
	return cfg, db
}

// logConfig 记录配置信息
func logConfig(cfg *config.Config) {
	log.Printf("应用名称: %s", cfg.App.Name)
	log.Printf("环境: %s", cfg.App.Environment)
	log.Printf("端口: %s", cfg.App.Port)
}

// startServer 启动服务器
func startServer(cfg *config.Config, db *gorm.DB) error {
	setGinMode(cfg.App.Environment)
	
	r := router.SetupRouter(db)
	log.Printf("服务器启动在端口 %s", cfg.App.Port)
	
	// 直接启动服务器（阻塞调用）
	return r.Run(":" + cfg.App.Port)
}

// setGinMode 设置Gin运行模式
func setGinMode(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

// closeDatabase 关闭数据库连接
func closeDatabase(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("数据库连接已关闭")
		}
	}
}