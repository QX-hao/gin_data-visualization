package config

import (
	"fmt"
)

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	Charset      string `yaml:"charset"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// Validate 验证数据库配置
func (c *DatabaseConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}

	if c.Port == "" {
		c.Port = "3306"
	}

	if c.Username == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}

	if c.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	if c.Charset == "" {
		c.Charset = "utf8mb4"
	}

	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 10
	}

	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 100
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset)
}

// GetDatabaseConfig 获取数据库配置（支持环境变量覆盖）
func GetDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnv("DB_PORT", "3306"),
		Username:     getEnv("DB_USERNAME", "root"),
		Password:     getEnv("DB_PASSWORD", "123456"),
		Database:     getEnv("DB_DATABASE", "gin_data_visualization"),
		Charset:      getEnv("DB_CHARSET", "utf8mb4"),
		MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 100),
	}
}

// GetInfo 获取数据库配置信息（隐藏密码）
func (c *DatabaseConfig) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"host":           c.Host,
		"port":           c.Port,
		"username":       c.Username,
		"database":       c.Database,
		"charset":        c.Charset,
		"max_idle_conns": c.MaxIdleConns,
		"max_open_conns": c.MaxOpenConns,
	}
}
