package repository

import (
	"errors"

	"github.com/adinegoro11/go-user-api/internal/model"
	"gorm.io/gorm"
)

var ErrProductNotFound = errors.New("product not found")

type ProductRepository interface {
	Create(product *model.Product) error
	FindByID(id uint) (*model.Product, error)
	FindByCode(code string) (*model.Product, error)
	FindAll(page int, pageSize int) ([]model.Product, int64, error)
	Update(product *model.Product) error
	Delete(id uint) error
}

type GormProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &GormProductRepository{db: db}
}

func (r *GormProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *GormProductRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	if err := r.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *GormProductRepository) FindByCode(code string) (*model.Product, error) {
	var product model.Product
	if err := r.db.Where("code = ?", code).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *GormProductRepository) FindAll(page int, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	if err := r.db.Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Order("id asc").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *GormProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *GormProductRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}
