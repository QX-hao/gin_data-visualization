package setting

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig 加载配置文件
func LoadConfig() error {
	// 设置配置文件路径 - 使用相对路径
	viper.SetConfigName("web-information")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../configs") // 相对于 pkg/setting 目录 -> go-web/configs
	viper.AddConfigPath("./configs")     // 相对于项目根目录 -> go-web/configs

	// 加载 web 配置
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("web-information.yaml 加载失败: %v", err)
	}

	// 加载数据库配置
	viper.SetConfigName("database")
	if err := viper.MergeInConfig(); err != nil {
		return fmt.Errorf("database.yaml 加载失败: %v", err)
	}

	// 加载日志配置
	viper.SetConfigName("logger")
	if err := viper.MergeInConfig(); err != nil {
		return fmt.Errorf("logger.yaml 加载失败: %v", err)
	}
	return nil
}
