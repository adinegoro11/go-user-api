package handler

import (
	"net/http"
	"strconv"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/repository"
	"github.com/adinegoro11/go-user-api/internal/service"
	"github.com/gin-gonic/gin"
)

type BillingService interface {
	Generate(customerID uint, subscriptionID uint, ottExternalID string) (dto.InvoiceResponse, error)
	FindByID(id uint) (dto.InvoiceResponse, error)
	FindByCustomerID(customerID uint) ([]dto.InvoiceResponse, error)
	MarkPaid(invoiceID uint) (dto.InvoiceResponse, error)
}

type BillingHandler struct {
	billingService BillingService
}

func NewBillingHandler(billingService BillingService) *BillingHandler {
	return &BillingHandler{billingService: billingService}
}

func (h *BillingHandler) Generate(c *gin.Context) {
	var req dto.GenerateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	invoice, err := h.billingService.Generate(req.CustomerID, req.SubscriptionID, "")
	if err != nil {
		switch err {
		case repository.ErrCustomerNotFound, repository.ErrSubscriptionNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		case service.ErrNoSubscriptionProducts:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate invoice"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "invoice generated", "data": invoice})
}

func (h *BillingHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	invoice, err := h.billingService.FindByID(uint(id))
	if err != nil {
		if err == repository.ErrInvoiceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "invoice not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get invoice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": invoice})
}

func (h *BillingHandler) ListByCustomer(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Query("customer_id"), 10, 64)
	if err != nil || customerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid customer_id"})
		return
	}

	invoices, err := h.billingService.FindByCustomerID(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get invoices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": invoices})
}

func (h *BillingHandler) Pay(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	invoice, err := h.billingService.MarkPaid(uint(id))
	if err != nil {
		if err == repository.ErrInvoiceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "invoice not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to pay invoice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "invoice paid", "data": invoice})
}
