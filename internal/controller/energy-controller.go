package controller

import (
	"go-web/internal/dao"
	"go-web/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetEnergyDistribution 获取能源类型分布
func GetEnergyDistribution(c *gin.Context) {
	// 从数据库获取数据
	energyList, err := dao.GetEnergyDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取数据失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    energyList,
	})
}
