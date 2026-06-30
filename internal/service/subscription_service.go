package service

import (
	"errors"
	"log/slog"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
)

var ErrInvalidSubscriptionInput = errors.New("invalid subscription input")
var ErrInvalidProductSelection = errors.New("invalid product selection")

type BillingGenerator interface {
	Generate(customerID uint, subscriptionID uint, ottExternalID string) (dto.InvoiceResponse, error)
}

type SubscriptionService struct {
	customerRepo     repository.CustomerRepository
	subscriptionRepo repository.SubscriptionRepository
	eventPublisher   EventPublisher
	notificationSvc  NotificationService
	ottSvc           OTTService
	billingGenerator BillingGenerator
}

func NewSubscriptionService(customerRepo repository.CustomerRepository, subscriptionRepo repository.SubscriptionRepository, eventPublisher EventPublisher, notificationSvc NotificationService, ottSvc OTTService, billingGenerator BillingGenerator) *SubscriptionService {
	return &SubscriptionService{
		customerRepo:     customerRepo,
		subscriptionRepo: subscriptionRepo,
		eventPublisher:   eventPublisher,
		notificationSvc:  notificationSvc,
		ottSvc:           ottSvc,
		billingGenerator: billingGenerator,
	}
}

func (s *SubscriptionService) Create(req dto.CreateSubscriptionRequest) (dto.SubscriptionResponse, error) {
	if req.CustomerID == 0 || len(req.ProductIDs) == 0 {
		return dto.SubscriptionResponse{}, ErrInvalidSubscriptionInput
	}

	customer, err := s.customerRepo.FindByID(req.CustomerID)
	if err != nil {
		return dto.SubscriptionResponse{}, err
	}

	productIDs := uniqueUint(req.ProductIDs)
	products, err := s.subscriptionRepo.FindProductsByIDs(productIDs)
	if err != nil {
		return dto.SubscriptionResponse{}, err
	}
	if len(products) != len(productIDs) {
		return dto.SubscriptionResponse{}, ErrInvalidProductSelection
	}

	subscription := model.Subscription{
		CustomerID: req.CustomerID,
		Status:     model.SubscriptionStatusWaitingInstallation,
	}
	if err := s.subscriptionRepo.Create(&subscription, productIDs); err != nil {
		return dto.SubscriptionResponse{}, err
	}

	_ = s.eventPublisher.Publish("SubscriptionCreated", map[string]interface{}{
		"subscription_id": subscription.ID,
		"customer_id":     req.CustomerID,
		"product_ids":     productIDs,
	})
	_ = s.eventPublisher.Publish("CustomerRegistered", map[string]interface{}{
		"customer_id":     customer.ID,
		"customer_number": customer.CustomerNumber,
		"subscription_id": subscription.ID,
	})

	_ = s.notificationSvc.SendWelcomeEmail(customer, products)
	ottID, _ := s.ottSvc.NotifyCustomer(customer)

	invoice, err := s.billingGenerator.Generate(req.CustomerID, subscription.ID, ottID)
	if err != nil {
		return dto.SubscriptionResponse{}, err
	}

	slog.Info("subscription created", "subscription_id", subscription.ID, "customer_id", req.CustomerID)
	return dto.SubscriptionResponse{
		ID:            subscription.ID,
		CustomerID:    subscription.CustomerID,
		Status:        subscription.Status,
		InvoiceID:     invoice.ID,
		InvoiceNumber: invoice.InvoiceNumber,
	}, nil
}

func uniqueUint(values []uint) []uint {
	seen := map[uint]struct{}{}
	result := make([]uint, 0, len(values))
	for _, v := range values {
		if v == 0 {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}
