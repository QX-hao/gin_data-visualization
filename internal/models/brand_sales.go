package models

// BrandSales 品牌销售数据模型
type BrandSales struct {
	BrandName  string  `json:"brand_name" gorm:"column:brand_name;type:varchar(255)"` // 品牌名称
	TotalSales float64 `json:"total_sales" gorm:"column:total_sales;type:decimal"`    // 总销售额
}

// TableName 指定表名
func (BrandSales) TableName() string {
	return "brand_sales_sum"
}
