package service

import (
	"errors"
	"testing"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"gorm.io/gorm"
)

type mockProductRepository struct {
	createFn     func(product *model.Product) error
	findByIDFn   func(id uint) (*model.Product, error)
	findByCodeFn func(code string) (*model.Product, error)
	findAllFn    func(page int, pageSize int) ([]model.Product, int64, error)
	updateFn     func(product *model.Product) error
	deleteFn     func(id uint) error
}

func (m *mockProductRepository) Create(product *model.Product) error {
	if m.createFn != nil {
		return m.createFn(product)
	}
	return nil
}

func (m *mockProductRepository) FindByID(id uint) (*model.Product, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, repository.ErrProductNotFound
}

func (m *mockProductRepository) FindByCode(code string) (*model.Product, error) {
	if m.findByCodeFn != nil {
		return m.findByCodeFn(code)
	}
	return nil, repository.ErrProductNotFound
}

func (m *mockProductRepository) FindAll(page int, pageSize int) ([]model.Product, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(page, pageSize)
	}
	return nil, 0, nil
}

func (m *mockProductRepository) Update(product *model.Product) error {
	if m.updateFn != nil {
		return m.updateFn(product)
	}
	return nil
}

func (m *mockProductRepository) Delete(id uint) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func TestProductServiceCreateSuccess(t *testing.T) {
	repo := &mockProductRepository{
		findByCodeFn: func(code string) (*model.Product, error) {
			if code != "INTERNET100" {
				t.Fatalf("unexpected normalized code: %s", code)
			}
			return nil, repository.ErrProductNotFound
		},
		createFn: func(product *model.Product) error {
			product.Model = gorm.Model{ID: 1}
			return nil
		},
	}

	svc := NewProductService(repo)
	resp, err := svc.Create(dto.CreateProductRequest{
		Code:        " internet100 ",
		Name:        " Internet 100 Mbps ",
		BillingType: "MONTHLY",
		Price:       300000,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ID != 1 {
		t.Fatalf("expected id 1, got %d", resp.ID)
	}
	if resp.Code != "INTERNET100" {
		t.Fatalf("expected normalized code, got %s", resp.Code)
	}
	if resp.BillingType != "monthly" {
		t.Fatalf("expected normalized billing type, got %s", resp.BillingType)
	}
}

func TestProductServiceCreateDuplicateCode(t *testing.T) {
	repo := &mockProductRepository{
		findByCodeFn: func(code string) (*model.Product, error) {
			return &model.Product{Model: gorm.Model{ID: 4}, Code: code}, nil
		},
	}

	svc := NewProductService(repo)
	_, err := svc.Create(dto.CreateProductRequest{
		Code:        "INTERNET100",
		Name:        "Internet 100 Mbps",
		BillingType: "monthly",
		Price:       300000,
	})
	if !errors.Is(err, ErrProductCodeAlreadyExists) {
		t.Fatalf("expected ErrProductCodeAlreadyExists, got %v", err)
	}
}

func TestProductServiceFindAllNormalizesPagination(t *testing.T) {
	repo := &mockProductRepository{
		findAllFn: func(page int, pageSize int) ([]model.Product, int64, error) {
			if page != 1 {
				t.Fatalf("expected normalized page 1, got %d", page)
			}
			if pageSize != 100 {
				t.Fatalf("expected normalized pageSize 100, got %d", pageSize)
			}
			return []model.Product{{Model: gorm.Model{ID: 1}, Code: "INTERNET100", Name: "Internet 100", BillingType: "monthly", Price: 300000, IsActive: true}}, 1, nil
		},
	}

	svc := NewProductService(repo)
	result, err := svc.FindAll(0, 500)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.TotalItems != 1 || result.TotalPages != 1 {
		t.Fatalf("unexpected totals: items=%d pages=%d", result.TotalItems, result.TotalPages)
	}
}

func TestProductServiceUpdateEmptyPayload(t *testing.T) {
	repo := &mockProductRepository{
		findByIDFn: func(id uint) (*model.Product, error) {
			return &model.Product{Model: gorm.Model{ID: id}, Code: "INTERNET100", Name: "Internet 100", BillingType: "monthly", Price: 300000, IsActive: true}, nil
		},
	}

	svc := NewProductService(repo)
	_, err := svc.Update(1, dto.UpdateProductRequest{})
	if !errors.Is(err, ErrInvalidProductInput) {
		t.Fatalf("expected ErrInvalidProductInput, got %v", err)
	}
}
