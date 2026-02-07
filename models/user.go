package models

import (
	"time"
)

type User struct {
	ID             int        `gorm:"primaryKey" json:"id"`
	EmployeeID     int        `gorm:"column:employee_id;not null;uniqueIndex" json:"employee_id"`
	PinHash        *string    `gorm:"column:pin_hash;size:255" json:"-"`
	FailedAttempts int        `gorm:"column:failed_attempts;default:0" json:"failed_attempts"`
	LockedUntil    *time.Time `gorm:"column:locked_until" json:"locked_until"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
