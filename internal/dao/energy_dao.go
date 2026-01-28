package dao

import (
    "go-web/internal/models"
)

// GetEnergyDistribution 获取能源类型分布数据
func GetEnergyDistribution() ([]models.EnergyType, error) {
    var energyList []models.EnergyType
    
    // 查询所有数据
    result := DB.Find(&energyList)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return energyList, nil
}