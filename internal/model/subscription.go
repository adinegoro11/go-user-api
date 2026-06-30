package model

import "gorm.io/gorm"

const (
	SubscriptionStatusWaitingInstallation = "WAITING_INSTALLATION"
	SubscriptionStatusActive              = "ACTIVE"
)

type Subscription struct {
	gorm.Model
	CustomerID uint   `gorm:"not null;index" json:"customer_id"`
	Status     string `gorm:"size:64;not null;default:WAITING_INSTALLATION" json:"status"`
}

type SubscriptionProduct struct {
	gorm.Model
	SubscriptionID uint `gorm:"not null;index" json:"subscription_id"`
	ProductID      uint `gorm:"not null;index" json:"product_id"`
}
