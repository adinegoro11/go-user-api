package repository

import (
	"github.com/adinegoro11/go-user-api/internal/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *model.DomainEvent) error
}

type EmailLogRepository interface {
	Create(log *model.EmailLog) error
}

type OTTRequestRepository interface {
	Create(req *model.OTTRequest) error
	Update(req *model.OTTRequest) error
}

type GormEventRepository struct {
	db *gorm.DB
}

type GormEmailLogRepository struct {
	db *gorm.DB
}

type GormOTTRequestRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &GormEventRepository{db: db}
}

func NewEmailLogRepository(db *gorm.DB) EmailLogRepository {
	return &GormEmailLogRepository{db: db}
}

func NewOTTRequestRepository(db *gorm.DB) OTTRequestRepository {
	return &GormOTTRequestRepository{db: db}
}

func (r *GormEventRepository) Create(event *model.DomainEvent) error {
	return r.db.Create(event).Error
}

func (r *GormEmailLogRepository) Create(log *model.EmailLog) error {
	return r.db.Create(log).Error
}

func (r *GormOTTRequestRepository) Create(req *model.OTTRequest) error {
	return r.db.Create(req).Error
}

func (r *GormOTTRequestRepository) Update(req *model.OTTRequest) error {
	return r.db.Save(req).Error
}
