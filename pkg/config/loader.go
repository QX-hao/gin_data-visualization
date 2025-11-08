package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构体
type Config struct {
	App      *AppConfig      `yaml:"app"`
	Database *DatabaseConfig `yaml:"database"`
	Server   *ServerConfig   `yaml:"server"`
	Log      *LogConfig      `yaml:"log"`
	JWT      *JWTConfig      `yaml:"jwt"`
}

var (
	globalConfig *Config
	configMutex  sync.RWMutex
	configPath   string = "config/app.yaml"
)

// LoadConfig 加载配置文件
func LoadConfig() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析YAML配置失败: %v", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	globalConfig = &config
	return nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()

	if globalConfig == nil {
		// 如果配置未加载，尝试加载默认配置
		return &Config{
			App:      GetAppConfig(),
			Database: GetDatabaseConfig(),
			Server:   GetServerConfig(),
			Log:      GetLogConfig(),
			JWT:      GetJWTConfig(),
		}
	}

	return globalConfig
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.App == nil {
		return fmt.Errorf("应用配置不能为空")
	}

	if err := c.App.Validate(); err != nil {
		return fmt.Errorf("应用配置验证失败: %v", err)
	}

	if c.Database == nil {
		return fmt.Errorf("数据库配置不能为空")
	}

	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("数据库配置验证失败: %v", err)
	}

	if c.Server == nil {
		return fmt.Errorf("服务器配置不能为空")
	}

	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("服务器配置验证失败: %v", err)
	}

	if c.Log == nil {
		return fmt.Errorf("日志配置不能为空")
	}

	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("日志配置验证失败: %v", err)
	}

	if c.JWT == nil {
		return fmt.Errorf("JWT配置不能为空")
	}

	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("JWT配置验证失败: %v", err)
	}

	return nil
}

// ReloadConfig 重新加载配置
func ReloadConfig() error {
	return LoadConfig()
}

// SetConfigPath 设置配置文件路径
func SetConfigPath(path string) {
	configMutex.Lock()
	defer configMutex.Unlock()
	configPath = path
}

// GetConfigInfo 获取配置信息（用于调试）
func (c *Config) GetConfigInfo() map[string]interface{} {
	return map[string]interface{}{
		"app":      c.App.GetInfo(),
		"database": c.Database.GetInfo(),
		"server":   c.Server.GetInfo(),
		"log":      c.Log.GetInfo(),
		"jwt":      c.JWT.GetInfo(),
	}
}