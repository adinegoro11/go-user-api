package model

import "gorm.io/gorm"

const (
	CustomerStatusPending = "PENDING"
	CustomerStatusActive  = "ACTIVE"
)

type Customer struct {
	gorm.Model
	CustomerNumber string `gorm:"size:32;uniqueIndex;not null" json:"customer_number"`
	Name           string `gorm:"size:255;not null" json:"name"`
	Email          string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Phone          string `gorm:"size:32;uniqueIndex;not null" json:"phone"`
	Address        string `gorm:"type:text;not null" json:"address"`
	Status         string `gorm:"size:32;not null;default:PENDING" json:"status"`
}
