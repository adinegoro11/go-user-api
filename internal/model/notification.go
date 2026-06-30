package model

import "gorm.io/gorm"

const (
	DeliveryStatusSent   = "SENT"
	DeliveryStatusFailed = "FAILED"
	DeliveryStatusQueued = "QUEUED"
)

type EmailLog struct {
	gorm.Model
	Recipient    string `gorm:"size:255;not null" json:"recipient"`
	Template     string `gorm:"size:64;not null" json:"template"`
	Status       string `gorm:"size:32;not null" json:"status"`
	ErrorMessage string `gorm:"type:text" json:"error_message"`
	Payload      string `gorm:"type:text" json:"payload"`
}

type OTTRequest struct {
	gorm.Model
	CustomerID      uint   `gorm:"not null;index" json:"customer_id"`
	Provider        string `gorm:"size:64;not null" json:"provider"`
	RequestPayload  string `gorm:"type:text;not null" json:"request_payload"`
	ResponsePayload string `gorm:"type:text" json:"response_payload"`
	Status          string `gorm:"size:32;not null" json:"status"`
	RetryCount      int    `gorm:"not null;default:0" json:"retry_count"`
}

type DomainEvent struct {
	gorm.Model
	Name    string `gorm:"size:128;not null;index" json:"name"`
	Payload string `gorm:"type:text;not null" json:"payload"`
	Status  string `gorm:"size:32;not null;default:PROCESSED" json:"status"`
}
