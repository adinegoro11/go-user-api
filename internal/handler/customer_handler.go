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

type CustomerService interface {
	Register(req dto.RegisterCustomerRequest) (dto.CustomerResponse, error)
	FindByID(id uint) (dto.CustomerResponse, error)
}

type CustomerHandler struct {
	customerService CustomerService
}

func NewCustomerHandler(customerService CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

func (h *CustomerHandler) Register(c *gin.Context) {
	var req dto.RegisterCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	customer, err := h.customerService.Register(req)
	if err != nil {
		switch err {
		case service.ErrInvalidCustomerInput:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case repository.ErrCustomerEmailAlreadyExists, repository.ErrCustomerPhoneAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		default:
			slog.Error("register customer failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to register customer"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "customer registered", "data": customer})
}

func (h *CustomerHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	customer, err := h.customerService.FindByID(uint(id))
	if err != nil {
		if err == repository.ErrCustomerNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}
