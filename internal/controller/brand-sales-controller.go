package controller

import (
	"go-web/internal/dao"
	"go-web/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBrandSales 获取所有品牌销售数据
func GetBrandSales(c *gin.Context) {
	// 从数据库获取数据
	brandSales, err := dao.GetBrandSales()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取品牌销售数据失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    brandSales,
	})
}

// GetTopBrandSales 获取前N名品牌销售数据
func GetTopBrandSales(c *gin.Context) {
	// 获取查询参数，默认为前20名
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // 默认显示前20名
	}

	// 限制最大数量，避免性能问题
	if limit > 100 {
		limit = 100
	}

	// 从数据库获取数据
	brandSales, err := dao.GetTopBrandSales(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取品牌销售数据失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    brandSales,
	})
}

// GetLastBrandSales 获取后N名品牌销售数据
func GetLastBrandSales(c *gin.Context) {
	// 获取查询参数，默认为后20名
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // 默认显示后20名
	}

	// 限制最大数量，避免性能问题
	if limit > 100 {
		limit = 100
	}

	// 从数据库获取数据
	brandSales, err := dao.GetLastBrandSales(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取品牌销售数据失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    brandSales,
	})
}
