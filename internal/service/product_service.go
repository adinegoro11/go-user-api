package service

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
)

var ErrProductCodeAlreadyExists = errors.New("product code already exists")
var ErrInvalidProductInput = errors.New("invalid product input")

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) Create(req dto.CreateProductRequest) (dto.ProductResponse, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	name := strings.TrimSpace(req.Name)
	billingType := strings.ToLower(strings.TrimSpace(req.BillingType))

	if code == "" || name == "" || !isValidBillingType(billingType) {
		slog.Warn("product create rejected due to invalid input", "code", code)
		return dto.ProductResponse{}, ErrInvalidProductInput
	}

	existing, err := s.productRepo.FindByCode(code)
	if err == nil && existing != nil {
		slog.Warn("product create rejected due to duplicate code", "code", code)
		return dto.ProductResponse{}, ErrProductCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, repository.ErrProductNotFound) {
		slog.Error("product create failed while checking code", "code", code, "error", err)
		return dto.ProductResponse{}, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	product := model.Product{
		Code:        code,
		Name:        name,
		BillingType: billingType,
		Price:       req.Price,
		IsActive:    isActive,
	}

	if err := s.productRepo.Create(&product); err != nil {
		slog.Error("product create failed", "code", code, "error", err)
		return dto.ProductResponse{}, err
	}

	slog.Info("product created", "product_id", product.ID, "code", product.Code)
	return toProductResponse(&product), nil
}

func (s *ProductService) FindByID(id uint) (dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return toProductResponse(product), nil
}

func (s *ProductService) FindAll(page int, pageSize int) (dto.PaginatedProductsResponse, error) {
	page = normalizePage(page)
	pageSize = normalizePageSize(pageSize)

	products, total, err := s.productRepo.FindAll(page, pageSize)
	if err != nil {
		slog.Error("product list failed", "error", err, "page", page, "page_size", pageSize)
		return dto.PaginatedProductsResponse{}, err
	}

	response := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		response = append(response, toProductResponse(&products[i]))
	}

	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}

	return dto.PaginatedProductsResponse{
		Data:       response,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func (s *ProductService) Update(id uint, req dto.UpdateProductRequest) (dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	if req.Code == nil && req.Name == nil && req.BillingType == nil && req.Price == nil && req.IsActive == nil {
		return dto.ProductResponse{}, ErrInvalidProductInput
	}

	if req.Code != nil {
		newCode := strings.ToUpper(strings.TrimSpace(*req.Code))
		if newCode == "" {
			return dto.ProductResponse{}, ErrInvalidProductInput
		}
		if newCode != product.Code {
			existing, err := s.productRepo.FindByCode(newCode)
			if err == nil && existing != nil {
				return dto.ProductResponse{}, ErrProductCodeAlreadyExists
			}
			if err != nil && !errors.Is(err, repository.ErrProductNotFound) {
				return dto.ProductResponse{}, err
			}
		}
		product.Code = newCode
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return dto.ProductResponse{}, ErrInvalidProductInput
		}
		product.Name = name
	}

	if req.BillingType != nil {
		billingType := strings.ToLower(strings.TrimSpace(*req.BillingType))
		if !isValidBillingType(billingType) {
			return dto.ProductResponse{}, ErrInvalidProductInput
		}
		product.BillingType = billingType
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.productRepo.Update(product); err != nil {
		slog.Error("product update failed", "product_id", id, "error", err)
		return dto.ProductResponse{}, err
	}

	slog.Info("product updated", "product_id", id, "code", product.Code)
	return toProductResponse(product), nil
}

func (s *ProductService) Delete(id uint) error {
	if err := s.productRepo.Delete(id); err != nil {
		slog.Error("product delete failed", "product_id", id, "error", err)
		return err
	}

	slog.Info("product soft deleted", "product_id", id)
	return nil
}

func toProductResponse(product *model.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          product.ID,
		Code:        product.Code,
		Name:        product.Name,
		BillingType: product.BillingType,
		Price:       product.Price,
		IsActive:    product.IsActive,
	}
}

func normalizePage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

func normalizePageSize(size int) int {
	if size < 1 {
		return 10
	}
	if size > 100 {
		return 100
	}
	return size
}

func isValidBillingType(billingType string) bool {
	return billingType == model.BillingTypeMonthly || billingType == model.BillingTypeOneTime
}
