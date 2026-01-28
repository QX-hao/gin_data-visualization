package dao

import (
	"go-web/internal/models"
)

// GetCitySales 获取城市销售数据
func GetCitySales() ([]models.CitySales, error) {
	var citySales []models.CitySales

	// 查询所有城市销售数据，按销售数量降序排列
	result := DB.Order("sales DESC").Find(&citySales)
	if result.Error != nil {
		return nil, result.Error
	}

	return citySales, nil
}

// GetTopCitySales 获取前N名城市销售数据
func GetTopCitySales(limit int) ([]models.CitySales, error) {
	var citySales []models.CitySales

	// 查询前N名城市销售数据，按销售数量降序排列
	result := DB.Order("sales DESC").Limit(limit).Find(&citySales)
	if result.Error != nil {
		return nil, result.Error
	}

	return citySales, nil
}
