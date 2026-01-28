package models

// EnergyType 能源类型分布模型
type EnergyType struct {
	EnergyName string `json:"energy_name" gorm:"column:energy_type;type:text"` // 能源类型名称字段
	Count      int64  `json:"count" gorm:"column:car_count;type:bigint"`       // 数量字段
}

// TableName 指定表名
func (EnergyType) TableName() string {
	return "ads_car_energy_distribution"
}
