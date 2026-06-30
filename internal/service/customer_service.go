package service

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
)

var ErrInvalidCustomerInput = errors.New("invalid customer input")

type CustomerService struct {
	customerRepo   repository.CustomerRepository
	eventPublisher EventPublisher
}

func NewCustomerService(customerRepo repository.CustomerRepository, eventPublisher EventPublisher) *CustomerService {
	return &CustomerService{customerRepo: customerRepo, eventPublisher: eventPublisher}
}

func (s *CustomerService) Register(req dto.RegisterCustomerRequest) (dto.CustomerResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	phone := strings.TrimSpace(req.Phone)
	address := strings.TrimSpace(req.Address)

	if name == "" || email == "" || phone == "" || address == "" {
		return dto.CustomerResponse{}, ErrInvalidCustomerInput
	}

	if existing, err := s.customerRepo.FindByEmail(email); err == nil && existing != nil {
		return dto.CustomerResponse{}, repository.ErrCustomerEmailAlreadyExists
	} else if err != nil && !errors.Is(err, repository.ErrCustomerNotFound) {
		return dto.CustomerResponse{}, err
	}

	if existing, err := s.customerRepo.FindByPhone(phone); err == nil && existing != nil {
		return dto.CustomerResponse{}, repository.ErrCustomerPhoneAlreadyExists
	} else if err != nil && !errors.Is(err, repository.ErrCustomerNotFound) {
		return dto.CustomerResponse{}, err
	}

	customerNumber, err := s.customerRepo.NextCustomerNumber()
	if err != nil {
		return dto.CustomerResponse{}, err
	}

	customer := model.Customer{
		CustomerNumber: customerNumber,
		Name:           name,
		Email:          email,
		Phone:          phone,
		Address:        address,
		Status:         model.CustomerStatusPending,
	}
	if err := s.customerRepo.Create(&customer); err != nil {
		return dto.CustomerResponse{}, err
	}

	_ = s.eventPublisher.Publish("CustomerRegistered", map[string]interface{}{
		"customer_id":     customer.ID,
		"customer_number": customer.CustomerNumber,
	})

	slog.Info("customer registered", "customer_id", customer.ID, "customer_number", customer.CustomerNumber)
	return dto.CustomerResponse{
		ID:             customer.ID,
		CustomerNumber: customer.CustomerNumber,
		Name:           customer.Name,
		Email:          customer.Email,
		Phone:          customer.Phone,
		Address:        customer.Address,
		Status:         customer.Status,
	}, nil
}

func (s *CustomerService) FindByID(id uint) (dto.CustomerResponse, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return dto.CustomerResponse{}, err
	}
	return dto.CustomerResponse{
		ID:             customer.ID,
		CustomerNumber: customer.CustomerNumber,
		Name:           customer.Name,
		Email:          customer.Email,
		Phone:          customer.Phone,
		Address:        customer.Address,
		Status:         customer.Status,
	}, nil
}
