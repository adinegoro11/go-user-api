package service

import (
	"encoding/json"
	"log/slog"

	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
)

type EventPublisher interface {
	Publish(name string, payload interface{}) error
}

type DBEventPublisher struct {
	eventRepo repository.EventRepository
}

func NewDBEventPublisher(eventRepo repository.EventRepository) *DBEventPublisher {
	return &DBEventPublisher{eventRepo: eventRepo}
}

func (p *DBEventPublisher) Publish(name string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	event := model.DomainEvent{
		Name:    name,
		Payload: string(payloadBytes),
		Status:  "PROCESSED",
	}
	if err := p.eventRepo.Create(&event); err != nil {
		return err
	}

	slog.Info("event published", "name", name, "event_id", event.ID)
	return nil
}
