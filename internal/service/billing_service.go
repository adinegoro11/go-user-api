package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
)

var ErrNoSubscriptionProducts = errors.New("subscription has no products")

type BillingService struct {
	invoiceRepo      repository.InvoiceRepository
	customerRepo     repository.CustomerRepository
	subscriptionRepo repository.SubscriptionRepository
	notificationSvc  NotificationService
}

func NewBillingService(invoiceRepo repository.InvoiceRepository, customerRepo repository.CustomerRepository, subscriptionRepo repository.SubscriptionRepository, notificationSvc NotificationService) *BillingService {
	return &BillingService{invoiceRepo: invoiceRepo, customerRepo: customerRepo, subscriptionRepo: subscriptionRepo, notificationSvc: notificationSvc}
}

func (s *BillingService) Generate(customerID uint, subscriptionID uint, ottExternalID string) (dto.InvoiceResponse, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		return dto.InvoiceResponse{}, err
	}

	subscription, err := s.subscriptionRepo.FindByID(subscriptionID)
	if err != nil {
		return dto.InvoiceResponse{}, err
	}
	if subscription.CustomerID != customerID {
		return dto.InvoiceResponse{}, repository.ErrSubscriptionNotFound
	}

	products, err := s.subscriptionRepo.FindProducts(subscriptionID)
	if err != nil {
		return dto.InvoiceResponse{}, err
	}
	if len(products) == 0 {
		return dto.InvoiceResponse{}, ErrNoSubscriptionProducts
	}

	items := make([]model.InvoiceItem, 0, len(products))
	var subtotal int64
	for _, p := range products {
		itemSubtotal := p.Price
		subtotal += itemSubtotal
		items = append(items, model.InvoiceItem{
			ProductName: p.Name,
			Quantity:    1,
			UnitPrice:   p.Price,
			Subtotal:    itemSubtotal,
			BillingType: p.BillingType,
			ProductCode: p.Code,
		})
	}

	tax := subtotal * 11 / 100
	discount := int64(0)
	if subtotal >= 1000000 {
		discount = subtotal * 5 / 100
	}
	grandTotal := subtotal + tax - discount

	invoiceNumber, err := s.invoiceRepo.NextInvoiceNumber(time.Now())
	if err != nil {
		return dto.InvoiceResponse{}, err
	}

	invoice := model.Invoice{
		InvoiceNumber:  invoiceNumber,
		CustomerID:     customerID,
		SubscriptionID: subscriptionID,
		Subtotal:       subtotal,
		Tax:            tax,
		Discount:       discount,
		GrandTotal:     grandTotal,
		DueDate:        time.Now().Add(7 * 24 * time.Hour),
		Status:         model.InvoiceStatusUnpaid,
		PaymentLink:    fmt.Sprintf("https://pay.local/invoices/%s", invoiceNumber),
		OTTExternalID:  ottExternalID,
	}

	if err := s.invoiceRepo.Create(&invoice, items); err != nil {
		return dto.InvoiceResponse{}, err
	}

	_ = s.notificationSvc.SendInvoiceEmail(customer, &invoice)
	slog.Info("invoice generated", "invoice_id", invoice.ID, "invoice_number", invoice.InvoiceNumber)

	return dto.InvoiceResponse{
		ID:             invoice.ID,
		InvoiceNumber:  invoice.InvoiceNumber,
		CustomerID:     invoice.CustomerID,
		SubscriptionID: invoice.SubscriptionID,
		Subtotal:       invoice.Subtotal,
		Tax:            invoice.Tax,
		Discount:       invoice.Discount,
		GrandTotal:     invoice.GrandTotal,
		Status:         invoice.Status,
		PaymentLink:    invoice.PaymentLink,
	}, nil
}

func (s *BillingService) FindByID(id uint) (dto.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return dto.InvoiceResponse{}, err
	}
	return dto.InvoiceResponse{
		ID:             invoice.ID,
		InvoiceNumber:  invoice.InvoiceNumber,
		CustomerID:     invoice.CustomerID,
		SubscriptionID: invoice.SubscriptionID,
		Subtotal:       invoice.Subtotal,
		Tax:            invoice.Tax,
		Discount:       invoice.Discount,
		GrandTotal:     invoice.GrandTotal,
		Status:         invoice.Status,
		PaymentLink:    invoice.PaymentLink,
	}, nil
}

func (s *BillingService) FindByCustomerID(customerID uint) ([]dto.InvoiceResponse, error) {
	invoices, err := s.invoiceRepo.FindByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	response := make([]dto.InvoiceResponse, 0, len(invoices))
	for _, invoice := range invoices {
		response = append(response, dto.InvoiceResponse{
			ID:             invoice.ID,
			InvoiceNumber:  invoice.InvoiceNumber,
			CustomerID:     invoice.CustomerID,
			SubscriptionID: invoice.SubscriptionID,
			Subtotal:       invoice.Subtotal,
			Tax:            invoice.Tax,
			Discount:       invoice.Discount,
			GrandTotal:     invoice.GrandTotal,
			Status:         invoice.Status,
			PaymentLink:    invoice.PaymentLink,
		})
	}
	return response, nil
}

func (s *BillingService) MarkPaid(invoiceID uint) (dto.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return dto.InvoiceResponse{}, err
	}

	invoice.Status = model.InvoiceStatusPaid
	if err := s.invoiceRepo.Update(invoice); err != nil {
		return dto.InvoiceResponse{}, err
	}

	subscription, err := s.subscriptionRepo.FindByID(invoice.SubscriptionID)
	if err == nil {
		subscription.Status = model.SubscriptionStatusActive
		_ = s.subscriptionRepo.Update(subscription)
	}

	customer, err := s.customerRepo.FindByID(invoice.CustomerID)
	if err == nil {
		customer.Status = model.CustomerStatusActive
		_ = s.customerRepo.Update(customer)
	}

	slog.Info("invoice paid", "invoice_id", invoice.ID)
	return dto.InvoiceResponse{
		ID:             invoice.ID,
		InvoiceNumber:  invoice.InvoiceNumber,
		CustomerID:     invoice.CustomerID,
		SubscriptionID: invoice.SubscriptionID,
		Subtotal:       invoice.Subtotal,
		Tax:            invoice.Tax,
		Discount:       invoice.Discount,
		GrandTotal:     invoice.GrandTotal,
		Status:         invoice.Status,
		PaymentLink:    invoice.PaymentLink,
	}, nil
}
