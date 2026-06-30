package service

import (
	"testing"
	"time"

	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"gorm.io/gorm"
)

type mockInvoiceRepo struct {
	createFn            func(invoice *model.Invoice, items []model.InvoiceItem) error
	findByIDFn          func(id uint) (*model.Invoice, error)
	findByCustomerIDFn  func(customerID uint) ([]model.Invoice, error)
	nextInvoiceNumberFn func(t time.Time) (string, error)
	updateFn            func(invoice *model.Invoice) error
}

func (m *mockInvoiceRepo) Create(invoice *model.Invoice, items []model.InvoiceItem) error {
	if m.createFn != nil {
		return m.createFn(invoice, items)
	}
	return nil
}
func (m *mockInvoiceRepo) FindByID(id uint) (*model.Invoice, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, repository.ErrInvoiceNotFound
}
func (m *mockInvoiceRepo) FindByCustomerID(customerID uint) ([]model.Invoice, error) {
	if m.findByCustomerIDFn != nil {
		return m.findByCustomerIDFn(customerID)
	}
	return nil, nil
}
func (m *mockInvoiceRepo) NextInvoiceNumber(t time.Time) (string, error) {
	if m.nextInvoiceNumberFn != nil {
		return m.nextInvoiceNumberFn(t)
	}
	return "INV-202607-000001", nil
}
func (m *mockInvoiceRepo) Update(invoice *model.Invoice) error {
	if m.updateFn != nil {
		return m.updateFn(invoice)
	}
	return nil
}

func TestBillingServiceGenerateSuccess(t *testing.T) {
	invoiceRepo := &mockInvoiceRepo{createFn: func(invoice *model.Invoice, items []model.InvoiceItem) error {
		invoice.Model = gorm.Model{ID: 1}
		return nil
	}}
	customerRepo := &mockCustomerRepo{findByIDFn: func(id uint) (*model.Customer, error) {
		return &model.Customer{Model: gorm.Model{ID: id}, Email: "john@email.com", Status: model.CustomerStatusPending}, nil
	}}
	subscriptionRepo := &mockSubscriptionRepo{
		findByIDFn: func(id uint) (*model.Subscription, error) {
			return &model.Subscription{Model: gorm.Model{ID: id}, CustomerID: 1, Status: model.SubscriptionStatusWaitingInstallation}, nil
		},
		findProductsFn: func(subscriptionID uint) ([]model.Product, error) {
			return []model.Product{{Model: gorm.Model{ID: 1}, Code: "INTERNET100", Name: "Internet", BillingType: "monthly", Price: 300000, IsActive: true}}, nil
		},
	}

	svc := NewBillingService(invoiceRepo, customerRepo, subscriptionRepo, &mockNotificationSvc{})
	resp, err := svc.Generate(1, 1, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != model.InvoiceStatusUnpaid {
		t.Fatalf("expected UNPAID, got %s", resp.Status)
	}
	if resp.GrandTotal != 333000 {
		t.Fatalf("expected grand total 333000, got %d", resp.GrandTotal)
	}
}

func TestBillingServiceMarkPaidActivatesSubscription(t *testing.T) {
	invoiceRepo := &mockInvoiceRepo{
		findByIDFn: func(id uint) (*model.Invoice, error) {
			return &model.Invoice{Model: model.Invoice{}.Model, InvoiceNumber: "INV-1", CustomerID: 1, SubscriptionID: 2, Status: model.InvoiceStatusUnpaid}, nil
		},
		updateFn: func(invoice *model.Invoice) error {
			return nil
		},
	}
	customerUpdated := false
	customerRepo := &mockCustomerRepo{
		findByIDFn: func(id uint) (*model.Customer, error) {
			return &model.Customer{Model: model.Customer{}.Model, Status: model.CustomerStatusPending}, nil
		},
		updateFn: func(customer *model.Customer) error {
			customerUpdated = customer.Status == model.CustomerStatusActive
			return nil
		},
	}
	subUpdated := false
	subscriptionRepo := &mockSubscriptionRepo{
		findByIDFn: func(id uint) (*model.Subscription, error) {
			return &model.Subscription{Model: model.Subscription{}.Model, Status: model.SubscriptionStatusWaitingInstallation}, nil
		},
		updateFn: func(subscription *model.Subscription) error {
			subUpdated = subscription.Status == model.SubscriptionStatusActive
			return nil
		},
	}

	svc := NewBillingService(invoiceRepo, customerRepo, subscriptionRepo, &mockNotificationSvc{})
	resp, err := svc.MarkPaid(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != model.InvoiceStatusPaid {
		t.Fatalf("expected PAID, got %s", resp.Status)
	}
	if !customerUpdated || !subUpdated {
		t.Fatalf("expected customer and subscription activation")
	}
}
