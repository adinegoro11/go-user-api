package handler

import (
	"net/http"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"github.com/adinegoro11/go-user-api/internal/service"
	"github.com/gin-gonic/gin"
)

type SubscriptionService interface {
	Create(req dto.CreateSubscriptionRequest) (dto.SubscriptionResponse, error)
}

type SubscriptionHandler struct {
	subscriptionService SubscriptionService
}

func NewSubscriptionHandler(subscriptionService SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionService: subscriptionService}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := h.subscriptionService.Create(req)
	if err != nil {
		switch err {
		case service.ErrInvalidSubscriptionInput, service.ErrInvalidProductSelection:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case repository.ErrCustomerNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "customer not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create subscription"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "subscription created", "data": result})
}
