package dao

import (
	"go-web/internal/models"
)

// GetCarLevelDistribution 获取所有汽车级别分布数据
func GetCarLevelDistribution() ([]models.CarLevelDistribution, error) {
	var carLevels []models.CarLevelDistribution

	// 查询所有汽车级别分布数据，按汽车数量降序排列
	result := DB.Order("car_count DESC").Find(&carLevels)
	if result.Error != nil {
		return nil, result.Error
	}

	return carLevels, nil
}

// GetTopCarLevels 获取前N名汽车级别分布数据
func GetTopCarLevels(limit int) ([]models.CarLevelDistribution, error) {
	var carLevels []models.CarLevelDistribution

	// 查询前N名汽车级别分布数据，按汽车数量降序排列
	result := DB.Order("car_count DESC").Limit(limit).Find(&carLevels)
	if result.Error != nil {
		return nil, result.Error
	}

	return carLevels, nil
}

// GetLastCarLevels 获取后N名汽车级别分布数据
func GetLastCarLevels(limit int) ([]models.CarLevelDistribution, error) {
	var carLevels []models.CarLevelDistribution

	// 查询后N名汽车级别分布数据，按汽车数量升序排列
	result := DB.Order("car_count ASC").Limit(limit).Find(&carLevels)
	if result.Error != nil {
		return nil, result.Error
	}

	return carLevels, nil
}
