package dao

import (
	"go-web/internal/models"
)

// GetBrandSales 获取所有品牌销售数据
func GetBrandSales() ([]models.BrandSales, error) {
	var brandSales []models.BrandSales

	// 查询所有品牌销售数据，按销售额降序排列
	result := DB.Order("total_sales DESC").Find(&brandSales)
	if result.Error != nil {
		return nil, result.Error
	}

	return brandSales, nil
}

// GetTopBrandSales 获取前N名品牌销售数据
func GetTopBrandSales(limit int) ([]models.BrandSales, error) {
	var brandSales []models.BrandSales

	// 查询前N名品牌销售数据，按销售额降序排列
	result := DB.Order("total_sales DESC").Limit(limit).Find(&brandSales)
	if result.Error != nil {
		return nil, result.Error
	}

	return brandSales, nil
}

// GetLastBrandSales 获取后N名品牌销售数据
func GetLastBrandSales(limit int) ([]models.BrandSales, error) {
	var brandSales []models.BrandSales

	// 查询后N名品牌销售数据，按销售额升序排列
	result := DB.Order("total_sales ASC").Limit(limit).Find(&brandSales)
	if result.Error != nil {
		return nil, result.Error
	}

	return brandSales, nil
}
