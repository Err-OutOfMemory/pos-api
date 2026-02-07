package models

type OrderDetail struct {
	ID             int     `gorm:"primaryKey;autoIncrement:false" json:"id"`
	OrderID        int     `gorm:"not null" json:"order_id"`
	ProductID      int     `gorm:"not null" json:"product_id"`
	DiscountAmount float64 `gorm:"column:discount_amount" json:"discount_amount"`
	Description    string  `gorm:"size:255" json:"description"`

	Order   Order   `gorm:"foreignKey:OrderID"`
	Product Product `gorm:"foreignKey:ProductID"`
}

func (OrderDetail) TableName() string {
	return "order_details"
}