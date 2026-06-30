package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
)

type OTTClient interface {
	NotifyCustomer(customer *model.Customer) (string, string, error)
}

type FakeOTTClient struct{}

func (c *FakeOTTClient) NotifyCustomer(customer *model.Customer) (string, string, error) {
	if os.Getenv("OTT_FORCE_FAIL") == "true" {
		return "", "{\"message\":\"forced failure\"}", errors.New("ott provider unavailable")
	}
	id := fmt.Sprintf("OTT-%06d", customer.ID)
	response := fmt.Sprintf("{\"id\":\"%s\",\"status\":\"ok\"}", id)
	return id, response, nil
}

type OTTService interface {
	NotifyCustomer(customer *model.Customer) (string, error)
}

type DefaultOTTService struct {
	ottRepo repository.OTTRequestRepository
	client  OTTClient
}

func NewDefaultOTTService(ottRepo repository.OTTRequestRepository, client OTTClient) *DefaultOTTService {
	if client == nil {
		client = &FakeOTTClient{}
	}
	return &DefaultOTTService{ottRepo: ottRepo, client: client}
}

func (s *DefaultOTTService) NotifyCustomer(customer *model.Customer) (string, error) {
	payloadBytes, _ := json.Marshal(map[string]interface{}{
		"customerNumber": customer.CustomerNumber,
		"name":           customer.Name,
		"email":          customer.Email,
	})

	ottID, responsePayload, err := s.client.NotifyCustomer(customer)
	status := model.DeliveryStatusSent
	retryCount := 0
	if err != nil {
		status = model.DeliveryStatusQueued
		retryCount = 1
	}

	ottReq := model.OTTRequest{
		CustomerID:      customer.ID,
		Provider:        "default-ott",
		RequestPayload:  string(payloadBytes),
		ResponsePayload: responsePayload,
		Status:          status,
		RetryCount:      retryCount,
	}
	if errCreate := s.ottRepo.Create(&ottReq); errCreate != nil {
		return "", errCreate
	}

	if err != nil {
		slog.Warn("ott notification queued for retry", "customer_id", customer.ID, "ott_request_id", ottReq.ID, "error", err)
		return "", nil
	}

	slog.Info("ott notification success", "customer_id", customer.ID, "ott_request_id", ottReq.ID)
	return ottID, nil
}
