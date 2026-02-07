package models

import "time"

type Order struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Date        time.Time `gorm:"type:date;not null" json:"date"`
	EmpID       int       `gorm:"column:emp_id;not null" json:"emp_id"`
	TotalPrice float64   `gorm:"not null" json:"total_price"`
	Discount   float64   `json:"discount"`
	OrderTypeID int      `gorm:"column:order_type_id;not null" json:"order_type_id"`

	User      User      `gorm:"foreignKey:EmpID"`
	OrderType OrderType `gorm:"foreignKey:OrderTypeID"`
}

func (Order) TableName() string {
	return "orders"
}
