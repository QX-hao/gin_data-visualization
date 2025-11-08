package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"gin_data-visualization/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitLogger 初始化日志系统
func InitLogger() {
	// 创建日志目录
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("创建日志目录失败: %v", err)
		return
	}

	// 设置日志文件路径
	logFile := filepath.Join(logDir, fmt.Sprintf("app_%s.log", time.Now().Format("20060102")))
	
	// 打开日志文件
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("打开日志文件失败: %v", err)
		return
	}

	// 设置日志输出到文件和控制台
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	
	log.Println("日志系统初始化完成")
}

// InitEnvironment 初始化环境变量
func InitEnvironment() {
	// 设置默认环境变量
	if os.Getenv("APP_ENV") == "" {
		os.Setenv("APP_ENV", "development")
	}
	
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
	
	log.Printf("环境变量初始化完成: APP_ENV=%s, PORT=%s", 
		os.Getenv("APP_ENV"), os.Getenv("PORT"))
}

// InitDatabaseWithRetry 带重试机制的数据库初始化
func InitDatabaseWithRetry(maxRetries int, retryInterval time.Duration) (*gorm.DB, error) {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Printf("配置文件加载失败，使用默认配置: %v", err)
	}
	
	cfg := config.GetConfig()
	
	for i := 0; i < maxRetries; i++ {
		db, err := initDatabase(cfg.Database)
		if err == nil {
			log.Printf("数据库连接成功 (第%d次重试)", i+1)
			return db, nil
		}
		
		log.Printf("数据库连接失败 (第%d次重试): %v", i+1, err)
		
		if i < maxRetries-1 {
			log.Printf("等待 %v 后重试...", retryInterval)
			time.Sleep(retryInterval)
		}
	}
	
	return nil, fmt.Errorf("数据库连接失败，已达到最大重试次数: %d", maxRetries)
}

// initDatabase 初始化数据库连接
func initDatabase(dbConfig *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := dbConfig.GetDSN()
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	log.Printf("数据库连接成功: %s@%s:%s/%s", 
		dbConfig.Username, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	
	return db, nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	// TODO: 添加具体的模型迁移
	log.Println("数据库迁移完成")
	return nil
}

// CheckSystemRequirements 检查系统要求
func CheckSystemRequirements() error {
	// 检查Go版本
	if runtime.Version() < "go1.18" {
		return fmt.Errorf("Go版本过低，需要1.18或更高版本，当前版本: %s", runtime.Version())
	}
	
	// 检查内存
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	// 检查磁盘空间（Windows兼容性处理）
	wd, _ := os.Getwd()
	info, err := os.Stat(wd)
	if err != nil {
		log.Printf("无法获取磁盘空间信息: %v", err)
	} else {
		// 在Windows上，我们无法直接获取磁盘空间，但可以检查目录是否可写
		if info.IsDir() {
			testFile := filepath.Join(wd, ".disk_space_test")
			file, err := os.Create(testFile)
			if err != nil {
				return fmt.Errorf("磁盘空间不足或目录不可写: %v", err)
			}
			file.Close()
			os.Remove(testFile)
		}
	}
	
	log.Println("系统要求检查通过")
	return nil
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	return map[string]interface{}{
		"go_version":     runtime.Version(),
		"go_os":          runtime.GOOS,
		"go_arch":        runtime.GOARCH,
		"num_cpu":        runtime.NumCPU(),
		"goroutines":     runtime.NumGoroutine(),
		"memory_alloc":   memStats.Alloc,
		"memory_total":   memStats.TotalAlloc,
		"memory_sys":     memStats.Sys,
		"memory_lookups": memStats.Lookups,
		"memory_mallocs": memStats.Mallocs,
		"memory_frees":   memStats.Frees,
		"start_time":     time.Now().Format("2006-01-02 15:04:05"),
	}
}

// ValidateConfig 验证配置
func ValidateConfig() error {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Printf("配置文件加载失败，使用默认配置: %v", err)
	}
	
	cfg := config.GetConfig()
	
	// 验证应用配置
	if cfg.App.Name == "" {
		return fmt.Errorf("应用名称不能为空")
	}
	
	if cfg.App.Port == "" {
		return fmt.Errorf("端口号不能为空")
	}
	
	// 验证数据库配置
	if cfg.Database.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}
	
	if cfg.Database.Username == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}
	
	if cfg.Database.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	
	log.Println("配置验证通过")
	return nil
}

// InitAll 初始化所有组件
func InitAll() (*gorm.DB, error) {
	// 检查系统要求
	if err := CheckSystemRequirements(); err != nil {
		return nil, err
	}
	
	// 初始化环境变量
	InitEnvironment()
	
	// 初始化日志系统
	InitLogger()
	
	// 验证配置
	if err := ValidateConfig(); err != nil {
		return nil, err
	}
	
	// 初始化数据库（带重试机制）
	db, err := InitDatabaseWithRetry(3, 5*time.Second)
	if err != nil {
		return nil, err
	}
	
	// 自动迁移数据库表
	if err := autoMigrate(db); err != nil {
		return nil, err
	}
	
	// 输出系统信息
	systemInfo := GetSystemInfo()
	log.Printf("系统信息: %+v", systemInfo)
	
	log.Println("所有组件初始化完成")
	return db, nil
}