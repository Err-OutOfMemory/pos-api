package models

type Category struct {
	ID int `gorm:"primaryKey" json:"id"`
	CategoryName string `gorm:"column:category_name;size:255;not null" json:"category_name"`
	Status       bool   `gorm:"not null" json:"status"`
}

func (Category) TableName() string {
	return "categories"
}
