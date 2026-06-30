package model

import "gorm.io/gorm"

const (
	BillingTypeMonthly = "monthly"
	BillingTypeOneTime = "one_time"
)

type Product struct {
	gorm.Model
	Code        string `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name        string `gorm:"size:255;not null" json:"name"`
	BillingType string `gorm:"size:32;not null" json:"billing_type"`
	Price       int64  `gorm:"not null" json:"price"`
	IsActive    bool   `gorm:"default:true;not null" json:"is_active"`
}
