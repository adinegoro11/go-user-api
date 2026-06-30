package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/adinegoro11/go-user-api/internal/model"
	"gorm.io/gorm"
)

var ErrCustomerNotFound = errors.New("customer not found")
var ErrCustomerEmailAlreadyExists = errors.New("customer email already exists")
var ErrCustomerPhoneAlreadyExists = errors.New("customer phone already exists")

type CustomerRepository interface {
	Create(customer *model.Customer) error
	FindByID(id uint) (*model.Customer, error)
	FindByEmail(email string) (*model.Customer, error)
	FindByPhone(phone string) (*model.Customer, error)
	NextCustomerNumber() (string, error)
	Update(customer *model.Customer) error
}

type GormCustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &GormCustomerRepository{db: db}
}

func (r *GormCustomerRepository) Create(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *GormCustomerRepository) FindByID(id uint) (*model.Customer, error) {
	var customer model.Customer
	if err := r.db.First(&customer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	return &customer, nil
}

func (r *GormCustomerRepository) FindByEmail(email string) (*model.Customer, error) {
	var customer model.Customer
	if err := r.db.Where("email = ?", email).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	return &customer, nil
}

func (r *GormCustomerRepository) FindByPhone(phone string) (*model.Customer, error) {
	var customer model.Customer
	if err := r.db.Where("phone = ?", phone).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	return &customer, nil
}

func (r *GormCustomerRepository) NextCustomerNumber() (string, error) {
	var count int64
	if err := r.db.Model(&model.Customer{}).Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("CUS%06d", count+1), nil
}

func (r *GormCustomerRepository) Update(customer *model.Customer) error {
	customer.UpdatedAt = time.Now()
	return r.db.Save(customer).Error
}
