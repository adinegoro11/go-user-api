package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"github.com/adinegoro11/go-user-api/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductService interface {
	Create(req dto.CreateProductRequest) (dto.ProductResponse, error)
	FindByID(id uint) (dto.ProductResponse, error)
	FindAll(page int, pageSize int) (dto.PaginatedProductsResponse, error)
	Update(id uint, req dto.UpdateProductRequest) (dto.ProductResponse, error)
	Delete(id uint) error
}

type ProductHandler struct {
	productService ProductService
}

func NewProductHandler(productService ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("product create validation failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	product, err := h.productService.Create(req)
	if err != nil {
		switch err {
		case service.ErrInvalidProductInput:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case service.ErrProductCodeAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		default:
			slog.Error("product create failed", "code", req.Code, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create product"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "create success", "product": product})
}

func (h *ProductHandler) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.productService.FindAll(page, pageSize)
	if err != nil {
		slog.Error("product list failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ProductHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	product, err := h.productService.FindByID(uint(id))
	if err != nil {
		if err == repository.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
			return
		}
		slog.Error("product detail failed", "product_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("product update validation failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	product, err := h.productService.Update(uint(id), req)
	if err != nil {
		switch err {
		case repository.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
		case service.ErrInvalidProductInput:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case service.ErrProductCodeAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		default:
			slog.Error("product update failed", "product_id", id, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update product"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update success", "product": product})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	if err := h.productService.Delete(uint(id)); err != nil {
		if err == repository.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
			return
		}
		slog.Error("product delete failed", "product_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete success"})
}
