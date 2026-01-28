package models

// CitySales 城市销售数据模型
type CitySales struct {
	City  string  `json:"city" gorm:"column:city;type:varchar(255)"` // 城市名称
	Sales float64 `json:"sales" gorm:"column:sales;type:decimal"`    // 销售数量
}

// TableName 指定表名
func (CitySales) TableName() string {
	return "citys_sv"
}
