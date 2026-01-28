package models

// CarLevelDistribution 汽车级别分布模型
type CarLevelDistribution struct {
	CarLevel string `json:"car_level" gorm:"column:car_level;type:text"`   // 汽车级别
	CarCount int64  `json:"car_count" gorm:"column:car_count;type:bigint"` // 汽车数量
}

// TableName 指定表名
func (CarLevelDistribution) TableName() string {
	return "ads_car_level_distribution"
}
