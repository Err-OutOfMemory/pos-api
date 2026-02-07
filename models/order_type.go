package models

type OrderType struct {
	ID     int    `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Type   string `gorm:"size:255;not null" json:"type"`
	Status string `gorm:"size:20;default:'active'" json:"status"`
}

func (OrderType) TableName() string {
	return "order_types"
}
