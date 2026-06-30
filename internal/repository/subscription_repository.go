package repository

import (
	"errors"

	"github.com/adinegoro11/go-user-api/internal/model"
	"gorm.io/gorm"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type SubscriptionRepository interface {
	Create(subscription *model.Subscription, productIDs []uint) error
	FindByID(id uint) (*model.Subscription, error)
	FindProducts(subscriptionID uint) ([]model.Product, error)
	FindProductsByIDs(productIDs []uint) ([]model.Product, error)
	Update(subscription *model.Subscription) error
}

type GormSubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &GormSubscriptionRepository{db: db}
}

func (r *GormSubscriptionRepository) Create(subscription *model.Subscription, productIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(subscription).Error; err != nil {
			return err
		}
		for _, productID := range productIDs {
			sp := model.SubscriptionProduct{SubscriptionID: subscription.ID, ProductID: productID}
			if err := tx.Create(&sp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormSubscriptionRepository) FindByID(id uint) (*model.Subscription, error) {
	var subscription model.Subscription
	if err := r.db.First(&subscription, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}
	return &subscription, nil
}

func (r *GormSubscriptionRepository) FindProducts(subscriptionID uint) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Table("products").
		Select("products.*").
		Joins("JOIN subscription_products ON subscription_products.product_id = products.id").
		Where("subscription_products.subscription_id = ?", subscriptionID).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *GormSubscriptionRepository) FindProductsByIDs(productIDs []uint) ([]model.Product, error) {
	var products []model.Product
	if err := r.db.Where("id IN ?", productIDs).Where("is_active = ?", true).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *GormSubscriptionRepository) Update(subscription *model.Subscription) error {
	return r.db.Save(subscription).Error
}
