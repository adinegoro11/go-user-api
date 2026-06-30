package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	InvoiceStatusUnpaid = "UNPAID"
	InvoiceStatusPaid   = "PAID"
	InvoiceStatusVoid   = "VOID"
)

type Invoice struct {
	gorm.Model
	InvoiceNumber  string    `gorm:"size:64;uniqueIndex;not null" json:"invoice_number"`
	CustomerID     uint      `gorm:"not null;index" json:"customer_id"`
	SubscriptionID uint      `gorm:"not null;index" json:"subscription_id"`
	Subtotal       int64     `gorm:"not null" json:"subtotal"`
	Tax            int64     `gorm:"not null" json:"tax"`
	Discount       int64     `gorm:"not null" json:"discount"`
	GrandTotal     int64     `gorm:"not null" json:"grand_total"`
	DueDate        time.Time `gorm:"not null" json:"due_date"`
	Status         string    `gorm:"size:32;not null;default:UNPAID" json:"status"`
	PaymentLink    string    `gorm:"size:255" json:"payment_link"`
	OTTExternalID  string    `gorm:"size:128" json:"ott_external_id"`
}

type InvoiceItem struct {
	gorm.Model
	InvoiceID   uint   `gorm:"not null;index" json:"invoice_id"`
	ProductName string `gorm:"size:255;not null" json:"product_name"`
	Quantity    int    `gorm:"not null" json:"quantity"`
	UnitPrice   int64  `gorm:"not null" json:"unit_price"`
	Subtotal    int64  `gorm:"not null" json:"subtotal"`
	BillingType string `gorm:"size:32;not null" json:"billing_type"`
	ProductCode string `gorm:"size:64;not null" json:"product_code"`
}
