package models

import "time"

type Order struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Date        time.Time `gorm:"type:date;not null" json:"date"`
	EmpID       int       `gorm:"column:emp_id;not null" json:"emp_id"`
	TotalPrice  float64   `gorm:"not null" json:"total_price"`
	Discount    float64   `json:"discount"`
	OrderTypeID int       `gorm:"column:order_type_id;not null" json:"order_type_id"`
	Status      string    `gorm:"size:20;default:'pending'" json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Employee     Employee      `gorm:"foreignKey:EmpID" json:"employee,omitzero"`
	OrderType    OrderType     `gorm:"foreignKey:OrderTypeID" json:"order_type,omitzero"`
	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID" json:"order_details,omitzero"`
}

func (Order) TableName() string {
	return "orders"
}
