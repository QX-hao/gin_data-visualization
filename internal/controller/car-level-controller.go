package controller

import (
	"go-web/internal/dao"
	"go-web/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCarLevelDistribution 获取所有汽车级别分布数据
func GetCarLevelDistribution(c *gin.Context) {
	// 从数据库获取数据
	carLevels, err := dao.GetCarLevelDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取汽车级别分布数据失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    carLevels,
	})
}

// GetTopCarLevels 获取前N名汽车级别分布数据
func GetTopCarLevels(c *gin.Context) {
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
	carLevels, err := dao.GetTopCarLevels(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取汽车级别分布数据失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    carLevels,
	})
}

// GetLastCarLevels 获取后N名汽车级别分布数据
func GetLastCarLevels(c *gin.Context) {
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
	carLevels, err := dao.GetLastCarLevels(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取汽车级别分布数据失败",
			Error:   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    carLevels,
	})
}
