package routers

import (
	"go-web/internal/controller"

	"github.com/gin-gonic/gin"
)

// SetupEnergyRoutes 设置能源相关路由
func SetupEnergyRoutes(api *gin.RouterGroup) {
	energy := api.Group("/energy")
	{
		energy.GET("/distribution", controller.GetEnergyDistribution)
		// 未来可以添加更多能源相关路由
		// energy.GET("/trend", controller.GetEnergyTrend)
		// energy.GET("/analysis", controller.GetEnergyAnalysis)
	}
}
