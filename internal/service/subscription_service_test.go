package service

import (
	"errors"
	"testing"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"gorm.io/gorm"
)

type mockSubscriptionRepo struct {
	createFn            func(subscription *model.Subscription, productIDs []uint) error
	findByIDFn          func(id uint) (*model.Subscription, error)
	findProductsFn      func(subscriptionID uint) ([]model.Product, error)
	findProductsByIDsFn func(productIDs []uint) ([]model.Product, error)
	updateFn            func(subscription *model.Subscription) error
}

func (m *mockSubscriptionRepo) Create(subscription *model.Subscription, productIDs []uint) error {
	if m.createFn != nil {
		return m.createFn(subscription, productIDs)
	}
	return nil
}
func (m *mockSubscriptionRepo) FindByID(id uint) (*model.Subscription, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, repository.ErrSubscriptionNotFound
}
func (m *mockSubscriptionRepo) FindProducts(subscriptionID uint) ([]model.Product, error) {
	if m.findProductsFn != nil {
		return m.findProductsFn(subscriptionID)
	}
	return nil, nil
}
func (m *mockSubscriptionRepo) FindProductsByIDs(productIDs []uint) ([]model.Product, error) {
	if m.findProductsByIDsFn != nil {
		return m.findProductsByIDsFn(productIDs)
	}
	return nil, nil
}
func (m *mockSubscriptionRepo) Update(subscription *model.Subscription) error {
	if m.updateFn != nil {
		return m.updateFn(subscription)
	}
	return nil
}

type mockNotificationSvc struct{}

func (m *mockNotificationSvc) SendWelcomeEmail(customer *model.Customer, products []model.Product) error {
	return nil
}
func (m *mockNotificationSvc) SendInvoiceEmail(customer *model.Customer, invoice *model.Invoice) error {
	return nil
}

type mockOTTSvc struct{}

func (m *mockOTTSvc) NotifyCustomer(customer *model.Customer) (string, error) {
	return "OTT-1", nil
}

type mockBillingGenerator struct {
	generateFn func(customerID uint, subscriptionID uint, ottExternalID string) (dto.InvoiceResponse, error)
}

func (m *mockBillingGenerator) Generate(customerID uint, subscriptionID uint, ottExternalID string) (dto.InvoiceResponse, error) {
	if m.generateFn != nil {
		return m.generateFn(customerID, subscriptionID, ottExternalID)
	}
	return dto.InvoiceResponse{}, nil
}

func TestSubscriptionServiceCreateSuccess(t *testing.T) {
	customerRepo := &mockCustomerRepo{findByIDFn: func(id uint) (*model.Customer, error) {
		return &model.Customer{Model: gorm.Model{ID: id}, CustomerNumber: "CUS000001", Email: "john@email.com", Name: "John"}, nil
	}}

	subRepo := &mockSubscriptionRepo{
		findProductsByIDsFn: func(productIDs []uint) ([]model.Product, error) {
			return []model.Product{{Model: gorm.Model{ID: 1}, Name: "Internet", Price: 100000, BillingType: "monthly", IsActive: true}}, nil
		},
		createFn: func(subscription *model.Subscription, productIDs []uint) error {
			subscription.ID = 10
			return nil
		},
	}

	svc := NewSubscriptionService(
		customerRepo,
		subRepo,
		&mockEventPublisher{},
		&mockNotificationSvc{},
		&mockOTTSvc{},
		&mockBillingGenerator{generateFn: func(customerID uint, subscriptionID uint, ottExternalID string) (dto.InvoiceResponse, error) {
			return dto.InvoiceResponse{ID: 22, InvoiceNumber: "INV-202607-000001"}, nil
		}},
	)

	resp, err := svc.Create(dto.CreateSubscriptionRequest{CustomerID: 1, ProductIDs: []uint{1}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ID != 10 || resp.InvoiceID != 22 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestSubscriptionServiceCreateInvalidProducts(t *testing.T) {
	customerRepo := &mockCustomerRepo{findByIDFn: func(id uint) (*model.Customer, error) {
		return &model.Customer{Model: gorm.Model{ID: id}}, nil
	}}
	subRepo := &mockSubscriptionRepo{findProductsByIDsFn: func(productIDs []uint) ([]model.Product, error) {
		return []model.Product{}, nil
	}}

	svc := NewSubscriptionService(customerRepo, subRepo, &mockEventPublisher{}, &mockNotificationSvc{}, &mockOTTSvc{}, &mockBillingGenerator{})
	_, err := svc.Create(dto.CreateSubscriptionRequest{CustomerID: 1, ProductIDs: []uint{100}})
	if !errors.Is(err, ErrInvalidProductSelection) {
		t.Fatalf("expected ErrInvalidProductSelection, got %v", err)
	}
}
