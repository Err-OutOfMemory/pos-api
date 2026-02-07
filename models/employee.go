package models

import (
	"time"
)

type Employee struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	EmpCode   string    `gorm:"column:emp_code;size:20;not null;uniqueIndex" json:"emp_code"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Role      string    `gorm:"size:100;not null" json:"role"`
	PhoneNum  string    `gorm:"column:phonenum;size:20" json:"phonenum"`
	Status    string    `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:EmployeeID" json:"user,omitzero"`
}

func (Employee) TableName() string {
	return "employees"
}
