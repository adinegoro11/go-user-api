package dto

type RegisterCustomerRequest struct {
	Name    string `json:"name" binding:"required,max=255"`
	Email   string `json:"email" binding:"required,email,max=255"`
	Phone   string `json:"phone" binding:"required,min=8,max=32"`
	Address string `json:"address" binding:"required,max=1000"`
}

type CustomerResponse struct {
	ID             uint   `json:"id"`
	CustomerNumber string `json:"customer_number"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	Status         string `json:"status"`
}

type CreateSubscriptionRequest struct {
	CustomerID uint   `json:"customer_id" binding:"required"`
	ProductIDs []uint `json:"product_ids" binding:"required,min=1,dive,gt=0"`
}

type SubscriptionResponse struct {
	ID            uint   `json:"id"`
	CustomerID    uint   `json:"customer_id"`
	Status        string `json:"status"`
	InvoiceID     uint   `json:"invoice_id,omitempty"`
	InvoiceNumber string `json:"invoice_number,omitempty"`
}

type GenerateInvoiceRequest struct {
	CustomerID     uint `json:"customer_id" binding:"required"`
	SubscriptionID uint `json:"subscription_id" binding:"required"`
}

type InvoiceResponse struct {
	ID             uint   `json:"id"`
	InvoiceNumber  string `json:"invoice_number"`
	CustomerID     uint   `json:"customer_id"`
	SubscriptionID uint   `json:"subscription_id"`
	Subtotal       int64  `json:"subtotal"`
	Tax            int64  `json:"tax"`
	Discount       int64  `json:"discount"`
	GrandTotal     int64  `json:"grand_total"`
	Status         string `json:"status"`
	PaymentLink    string `json:"payment_link"`
}
