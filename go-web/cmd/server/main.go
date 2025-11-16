package main

import (
	"go-web/internal/routers"
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	// 使用viper加载 go-web\configs\web-information.yaml 文件
	viper.SetConfigName("web-information")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.ReadInConfig()

	engine := routers.SetupRouter()
	engine.Run(fmt.Sprintf("%s:%d", viper.GetString("AppName.Host"), viper.GetInt("AppName.Port")))
}