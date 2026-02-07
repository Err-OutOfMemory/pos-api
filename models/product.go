package models

type Product struct {
	ID          int     `gorm:"primaryKey" json:"id"`
	ProductName string  `gorm:"column:product_name;size:255;not null" json:"product_name"`
	CategoryID  int     `gorm:"not null" json:"category_id"`
	Description string  `gorm:"size:255" json:"description"`
	Type        string  `gorm:"size:255;not null;default:'unit'" json:"type"`
	Price       float64 `gorm:"not null" json:"price"`
	ImgPath     string  `gorm:"column:img_path;size:255" json:"img_path"`
	Status      string  `gorm:"size:20;default:'active'" json:"status"`

	Category    Category `gorm:"foreignKey:CategoryID" json:"category,omitzero"`
}

func (Product) TableName() string {
	return "products"
}
