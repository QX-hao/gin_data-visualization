package config

import (
	"fmt"
	"os"
	"strconv"
)

// AppConfig 应用配置结构体
type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
	Port        string `yaml:"port"`
	Debug       bool   `yaml:"debug"`
}

// LogConfig 日志配置结构体
type LogConfig struct {
	Level      string `yaml:"level"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// JWTConfig JWT配置结构体
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireTime int64  `yaml:"expire_time"`
}

// Validate 验证应用配置
func (c *AppConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("应用名称不能为空")
	}

	if c.Version == "" {
		return fmt.Errorf("应用版本不能为空")
	}

	if c.Environment == "" {
		c.Environment = "development"
	}

	if c.Port == "" {
		c.Port = "8080"
	}

	return nil
}

// Validate 验证日志配置
func (c *LogConfig) Validate() error {
	if c.Level == "" {
		c.Level = "info"
	}

	if c.FilePath == "" {
		c.FilePath = "./logs/app.log"
	}

	if c.MaxSize == 0 {
		c.MaxSize = 10
	}

	if c.MaxBackups == 0 {
		c.MaxBackups = 5
	}

	if c.MaxAge == 0 {
		c.MaxAge = 30
	}

	return nil
}

// GetInfo 获取日志配置信息
func (c *LogConfig) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"level":       c.Level,
		"file_path":   c.FilePath,
		"max_size":    c.MaxSize,
		"max_backups": c.MaxBackups,
		"max_age":     c.MaxAge,
		"compress":    c.Compress,
	}
}

// Validate 验证JWT配置
func (c *JWTConfig) Validate() error {
	if c.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	if c.ExpireTime == 0 {
		c.ExpireTime = 86400
	}

	return nil
}

// GetInfo 获取JWT配置信息
func (c *JWTConfig) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"secret":      c.Secret,
		"expire_time": c.ExpireTime,
	}
}

// IsDevelopment 判断是否为开发环境
func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction 判断是否为生产环境
func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

// IsStaging 判断是否为测试环境
func (c *AppConfig) IsStaging() bool {
	return c.Environment == "staging"
}

// GetServerAddress 获取服务器地址
func (c *AppConfig) GetServerAddress() string {
	return ":" + c.Port
}

// GetInfo 获取应用信息
func (c *AppConfig) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":        c.Name,
		"version":     c.Version,
		"environment": c.Environment,
		"port":        c.Port,
		"debug":       c.Debug,
	}
}

// GetLogConfig 获取日志配置（支持环境变量覆盖）
func GetLogConfig() *LogConfig {
	return &LogConfig{
		Level:      getEnv("LOG_LEVEL", "info"),
		FilePath:   getEnv("LOG_FILE_PATH", "./logs/app.log"),
		MaxSize:    getEnvInt("LOG_MAX_SIZE", 10),
		MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 5),
		MaxAge:     getEnvInt("LOG_MAX_AGE", 30),
		Compress:   getEnvBool("LOG_COMPRESS", true),
	}
}

// GetAppConfig 获取应用配置（支持环境变量覆盖）
func GetAppConfig() *AppConfig {
	return &AppConfig{
		Name:        getEnv("APP_NAME", "Gin Data Visualization"),
		Version:     getEnv("APP_VERSION", "1.0.0"),
		Environment: getEnv("APP_ENVIRONMENT", "development"),
		Port:        getEnv("APP_PORT", "8080"),
		Debug:       getEnvBool("APP_DEBUG", true),
	}
}

// GetJWTConfig 获取JWT配置（支持环境变量覆盖）
func GetJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret:     getEnv("JWT_SECRET", "your-secret-key"),
		ExpireTime: getEnvInt64("JWT_EXPIRE_TIME", 86400),
	}
}

// 辅助函数：获取环境变量
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return intValue
}