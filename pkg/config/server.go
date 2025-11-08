package config

import (
	"time"
)

// ServerConfig 服务器配置结构体
type ServerConfig struct {
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// Validate 验证服务器配置
func (c *ServerConfig) Validate() error {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 10 * time.Second
	}

	if c.WriteTimeout == 0 {
		c.WriteTimeout = 10 * time.Second
	}

	if c.IdleTimeout == 0 {
		c.IdleTimeout = 60 * time.Second
	}

	return nil
}

// GetServerConfig 获取服务器配置（支持环境变量覆盖）
func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		ReadTimeout:  parseDuration(getEnv("SERVER_READ_TIMEOUT", "10s"), 10*time.Second),
		WriteTimeout: parseDuration(getEnv("SERVER_WRITE_TIMEOUT", "10s"), 10*time.Second),
		IdleTimeout:  parseDuration(getEnv("SERVER_IDLE_TIMEOUT", "60s"), 60*time.Second),
	}
}

// parseDuration 解析时间字符串，如果解析失败则返回默认值
func parseDuration(value string, defaultValue time.Duration) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}

// GetInfo 获取服务器配置信息
func (c *ServerConfig) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"read_timeout":  c.ReadTimeout.String(),
		"write_timeout": c.WriteTimeout.String(),
		"idle_timeout":  c.IdleTimeout.String(),
	}
}