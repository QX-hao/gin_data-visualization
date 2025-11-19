package dao

import (
	"fmt"
	"time"
	"go-web/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// InitDB 初始化数据库连接
func InitDB() error {
	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"))

	logger.Infow("正在连接数据库",
		"user", viper.GetString("mysql.user"),
		"host", viper.GetString("mysql.host"),
		"port", viper.GetInt("mysql.port"),
		"database", viper.GetString("mysql.dbname"))

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用默认的事务（为了保持一致性，建议在业务层控制事务）
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logger.Errorw("数据库连接失败", "error", err)
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 获取通用数据库对象 sql.DB，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorw("获取数据库连接池失败", "error", err)
		return fmt.Errorf("获取数据库连接池失败: %v", err)
	}

	// 设置数据库连接池参数
	sqlDB.SetMaxIdleConns(10)                  // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)                 // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour)        // 设置了连接可复用的最大时间

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		logger.Errorw("数据库连接测试失败", "error", err)
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	DB = db
	logger.Info("数据库连接成功")
	
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			logger.Errorw("获取数据库连接失败", "error", err)
			return err
		}
		
		if err := sqlDB.Close(); err != nil {
			logger.Errorw("关闭数据库连接失败", "error", err)
			return err
		}
		
		logger.Info("数据库连接已关闭")
	}
	return nil
}

