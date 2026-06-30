package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/middleware"
	"github.com/adinegoro11/go-user-api/internal/service"
	"github.com/gin-gonic/gin"
)

type mockProductService struct {
	createFn   func(req dto.CreateProductRequest) (dto.ProductResponse, error)
	findByIDFn func(id uint) (dto.ProductResponse, error)
	findAllFn  func(page int, pageSize int) (dto.PaginatedProductsResponse, error)
	updateFn   func(id uint, req dto.UpdateProductRequest) (dto.ProductResponse, error)
	deleteFn   func(id uint) error
}

func (m *mockProductService) Create(req dto.CreateProductRequest) (dto.ProductResponse, error) {
	if m.createFn != nil {
		return m.createFn(req)
	}
	return dto.ProductResponse{}, nil
}

func (m *mockProductService) FindByID(id uint) (dto.ProductResponse, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return dto.ProductResponse{}, nil
}

func (m *mockProductService) FindAll(page int, pageSize int) (dto.PaginatedProductsResponse, error) {
	if m.findAllFn != nil {
		return m.findAllFn(page, pageSize)
	}
	return dto.PaginatedProductsResponse{}, nil
}

func (m *mockProductService) Update(id uint, req dto.UpdateProductRequest) (dto.ProductResponse, error) {
	if m.updateFn != nil {
		return m.updateFn(id, req)
	}
	return dto.ProductResponse{}, nil
}

func (m *mockProductService) Delete(id uint) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func TestProductHandlerCreateDuplicateCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewProductHandler(&mockProductService{
		createFn: func(req dto.CreateProductRequest) (dto.ProductResponse, error) {
			return dto.ProductResponse{}, service.ErrProductCodeAlreadyExists
		},
	})
	r.POST("/products", h.Create)

	body, _ := json.Marshal(gin.H{"code": "INTERNET100", "name": "Internet 100", "billing_type": "monthly", "price": 1000})
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestProductHandlerCreateInvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewProductHandler(&mockProductService{})
	r.POST("/products", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("{bad-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestProductHandlerCreateServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewProductHandler(&mockProductService{
		createFn: func(req dto.CreateProductRequest) (dto.ProductResponse, error) {
			return dto.ProductResponse{}, errors.New("database down")
		},
	})
	r.POST("/products", h.Create)

	body, _ := json.Marshal(gin.H{"code": "INTERNET100", "name": "Internet 100", "billing_type": "monthly", "price": 1000})
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestProductCreateRouteAdminOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewProductHandler(&mockProductService{})

	r.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})
	products := r.Group("/products")
	products.Use(middleware.RoleMiddleware("admin"))
	products.POST("", h.Create)

	body, _ := json.Marshal(gin.H{"code": "INTERNET100", "name": "Internet 100", "billing_type": "monthly", "price": 1000})
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestProductUpdateRouteAdminOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewProductHandler(&mockProductService{})

	r.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})
	products := r.Group("/products")
	products.Use(middleware.RoleMiddleware("admin"))
	products.PUT("/:id", h.Update)

	body, _ := json.Marshal(gin.H{"name": "Updated Name"})
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestProductDeleteRouteAdminOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewProductHandler(&mockProductService{})

	r.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})
	products := r.Group("/products")
	products.Use(middleware.RoleMiddleware("admin"))
	products.DELETE("/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}
