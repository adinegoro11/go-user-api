package service

import (
	"errors"
	"testing"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"gorm.io/gorm"
)

type mockCustomerRepo struct {
	createFn             func(customer *model.Customer) error
	findByIDFn           func(id uint) (*model.Customer, error)
	findByEmailFn        func(email string) (*model.Customer, error)
	findByPhoneFn        func(phone string) (*model.Customer, error)
	nextCustomerNumberFn func() (string, error)
	updateFn             func(customer *model.Customer) error
}

func (m *mockCustomerRepo) Create(customer *model.Customer) error {
	if m.createFn != nil {
		return m.createFn(customer)
	}
	return nil
}

func (m *mockCustomerRepo) FindByID(id uint) (*model.Customer, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, repository.ErrCustomerNotFound
}

func (m *mockCustomerRepo) FindByEmail(email string) (*model.Customer, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(email)
	}
	return nil, repository.ErrCustomerNotFound
}

func (m *mockCustomerRepo) FindByPhone(phone string) (*model.Customer, error) {
	if m.findByPhoneFn != nil {
		return m.findByPhoneFn(phone)
	}
	return nil, repository.ErrCustomerNotFound
}

func (m *mockCustomerRepo) NextCustomerNumber() (string, error) {
	if m.nextCustomerNumberFn != nil {
		return m.nextCustomerNumberFn()
	}
	return "CUS000001", nil
}

func (m *mockCustomerRepo) Update(customer *model.Customer) error {
	if m.updateFn != nil {
		return m.updateFn(customer)
	}
	return nil
}

type mockEventPublisher struct {
	publishFn func(name string, payload interface{}) error
}

func (m *mockEventPublisher) Publish(name string, payload interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(name, payload)
	}
	return nil
}

func TestCustomerServiceRegisterSuccess(t *testing.T) {
	repo := &mockCustomerRepo{
		findByEmailFn: func(email string) (*model.Customer, error) {
			return nil, repository.ErrCustomerNotFound
		},
		findByPhoneFn: func(phone string) (*model.Customer, error) {
			return nil, repository.ErrCustomerNotFound
		},
		createFn: func(customer *model.Customer) error {
			customer.ID = 1
			return nil
		},
	}

	published := false
	svc := NewCustomerService(repo, &mockEventPublisher{publishFn: func(name string, payload interface{}) error {
		published = name == "CustomerRegistered"
		return nil
	}})

	resp, err := svc.Register(dto.RegisterCustomerRequest{Name: "John", Email: "JOHN@EMAIL.COM", Phone: "0812", Address: "Jakarta"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.ID != 1 {
		t.Fatalf("expected id 1, got %d", resp.ID)
	}
	if resp.Email != "john@email.com" {
		t.Fatalf("expected normalized email, got %s", resp.Email)
	}
	if !published {
		t.Fatalf("expected customer event published")
	}
}

func TestCustomerServiceRegisterDuplicateEmail(t *testing.T) {
	repo := &mockCustomerRepo{
		findByEmailFn: func(email string) (*model.Customer, error) {
			return &model.Customer{Model: gorm.Model{ID: 2}, Email: email}, nil
		},
	}
	svc := NewCustomerService(repo, &mockEventPublisher{})

	_, err := svc.Register(dto.RegisterCustomerRequest{Name: "John", Email: "john@email.com", Phone: "0812", Address: "Jakarta"})
	if !errors.Is(err, repository.ErrCustomerEmailAlreadyExists) {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
}
