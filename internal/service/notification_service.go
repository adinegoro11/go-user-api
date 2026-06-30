package service

import (
	"encoding/json"
	"log/slog"

	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
)

type NotificationService interface {
	SendWelcomeEmail(customer *model.Customer, products []model.Product) error
	SendInvoiceEmail(customer *model.Customer, invoice *model.Invoice) error
}

type EmailNotificationService struct {
	emailLogRepo repository.EmailLogRepository
}

func NewEmailNotificationService(emailLogRepo repository.EmailLogRepository) *EmailNotificationService {
	return &EmailNotificationService{emailLogRepo: emailLogRepo}
}

func (s *EmailNotificationService) SendWelcomeEmail(customer *model.Customer, products []model.Product) error {
	names := make([]string, 0, len(products))
	for _, p := range products {
		names = append(names, p.Name)
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"customer_number":  customer.CustomerNumber,
		"customer_name":    customer.Name,
		"selected_package": names,
	})

	emailLog := model.EmailLog{
		Recipient: customer.Email,
		Template:  "welcome",
		Status:    model.DeliveryStatusSent,
		Payload:   string(payload),
	}
	if err := s.emailLogRepo.Create(&emailLog); err != nil {
		return err
	}

	slog.Info("welcome email logged", "customer_id", customer.ID, "email_log_id", emailLog.ID)
	return nil
}

func (s *EmailNotificationService) SendInvoiceEmail(customer *model.Customer, invoice *model.Invoice) error {
	payload, _ := json.Marshal(map[string]interface{}{
		"invoice_number": invoice.InvoiceNumber,
		"amount":         invoice.GrandTotal,
		"payment_link":   invoice.PaymentLink,
	})

	emailLog := model.EmailLog{
		Recipient: customer.Email,
		Template:  "invoice",
		Status:    model.DeliveryStatusSent,
		Payload:   string(payload),
	}
	if err := s.emailLogRepo.Create(&emailLog); err != nil {
		return err
	}

	slog.Info("invoice email logged", "customer_id", customer.ID, "invoice_id", invoice.ID, "email_log_id", emailLog.ID)
	return nil
}
