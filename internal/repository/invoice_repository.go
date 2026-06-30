package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/adinegoro11/go-user-api/internal/model"
	"gorm.io/gorm"
)

var ErrInvoiceNotFound = errors.New("invoice not found")

type InvoiceRepository interface {
	Create(invoice *model.Invoice, items []model.InvoiceItem) error
	FindByID(id uint) (*model.Invoice, error)
	FindByCustomerID(customerID uint) ([]model.Invoice, error)
	NextInvoiceNumber(t time.Time) (string, error)
	Update(invoice *model.Invoice) error
}

type GormInvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &GormInvoiceRepository{db: db}
}

func (r *GormInvoiceRepository) Create(invoice *model.Invoice, items []model.InvoiceItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(invoice).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].InvoiceID = invoice.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormInvoiceRepository) FindByID(id uint) (*model.Invoice, error) {
	var invoice model.Invoice
	if err := r.db.First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvoiceNotFound
		}
		return nil, err
	}
	return &invoice, nil
}

func (r *GormInvoiceRepository) FindByCustomerID(customerID uint) ([]model.Invoice, error) {
	var invoices []model.Invoice
	if err := r.db.Where("customer_id = ?", customerID).Order("id desc").Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

func (r *GormInvoiceRepository) NextInvoiceNumber(t time.Time) (string, error) {
	prefix := t.Format("200601")
	var count int64
	if err := r.db.Model(&model.Invoice{}).Where("to_char(created_at, 'YYYYMM') = ?", prefix).Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("INV-%s-%06d", prefix, count+1), nil
}

func (r *GormInvoiceRepository) Update(invoice *model.Invoice) error {
	return r.db.Save(invoice).Error
}
