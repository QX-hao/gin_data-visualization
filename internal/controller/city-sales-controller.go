package controller

import (
	"go-web/internal/dao"
	"go-web/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCitySales 获取所有城市销售数据
func GetCitySales(c *gin.Context) {
	// 从数据库获取数据
	citySales, err := dao.GetCitySales()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取城市销售数据失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    citySales,
	})
}

// GetTopCitySales 获取前N名城市销售数据
func GetTopCitySales(c *gin.Context) {
	// 获取查询参数，默认为前10名
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // 默认显示前10名
	}

	// 限制最大数量，避免性能问题
	if limit > 100 {
		limit = 100
	}

	// 从数据库获取数据
	citySales, err := dao.GetTopCitySales(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取城市销售数据失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    citySales,
	})
}
