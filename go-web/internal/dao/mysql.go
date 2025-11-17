package dao

import (
	"fmt"
	"time"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	logger *zap.Logger
)

// InitLogger 初始化日志记录器
func InitLogger() error {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		return fmt.Errorf("初始化日志记录器失败: %v", err)
	}
	return nil
}

// InitDB 初始化数据库连接
func InitDB() error {
	// 初始化日志记录器
	if err := InitLogger(); err != nil {
		return err
	}
	defer logger.Sync()

	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"))

	logger.Info("正在连接数据库",
		zap.String("user", viper.GetString("mysql.user")),
		zap.String("host", viper.GetString("mysql.host")),
		zap.Int("port", viper.GetInt("mysql.port")),
		zap.String("database", viper.GetString("mysql.dbname")))

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用默认的事务（为了保持一致性，建议在业务层控制事务）
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logger.Error("数据库连接失败", zap.Error(err))
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 获取通用数据库对象 sql.DB，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("获取数据库连接池失败", zap.Error(err))
		return fmt.Errorf("获取数据库连接池失败: %v", err)
	}

	// 设置数据库连接池参数
	sqlDB.SetMaxIdleConns(10)                  // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)                 // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour)        // 设置了连接可复用的最大时间

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		logger.Error("数据库连接测试失败", zap.Error(err))
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

// GetLogger 获取日志记录器实例
func GetLogger() *zap.Logger {
	return logger
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			logger.Error("获取数据库连接失败", zap.Error(err))
			return err
		}
		
		if err := sqlDB.Close(); err != nil {
			logger.Error("关闭数据库连接失败", zap.Error(err))
			return err
		}
		
		logger.Info("数据库连接已关闭")
	}
	return nil
}

