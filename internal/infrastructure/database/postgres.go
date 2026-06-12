package database

import (
	"context"
	"errors"
	"fmt"
	"govent/internal/domain/shared"
	"govent/internal/domain/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventPostgresRepository struct {
	db *gorm.DB
}

func NewEventPostgresRepository(db *gorm.DB) types.EventRepository {
	return &EventPostgresRepository{db: db}
}

func (r *EventPostgresRepository) Create(ctx context.Context, event *types.Event) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(event).Error; err != nil {
			return fmt.Errorf("error inserting main event: %w", err)
		}

		var subscribers []string
		err := tx.Table("subscriptions").
			Where("event_name = ? AND source = ?", event.Name, event.Source).
			Pluck("subscriber_name", &subscribers).Error

		if err != nil {
			return fmt.Errorf("error searching subscriptions: %w", err)
		}

		for _, subscriber := range subscribers {
			queueId := uuid.New().String()

			msg, err := types.NewQueueMessage(queueId, subscriber, event.Id)
			if err != nil {
				return fmt.Errorf("error al construir QueueMessage: %w", err)
			}

			if err := tx.Table("queue").Create(msg).Error; err != nil {
				return fmt.Errorf("error al insertar en la cola de %s: %w", subscriber, err)
			}
		}

		return nil
	})
}

func (r *EventPostgresRepository) One(ctx context.Context, id *types.SharedId) (*types.Event, error) {
	var event types.Event

	err := r.db.WithContext(ctx).
		First(&event, "id = ?", id.Value()).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrEventNotFound
		}
		return nil, shared.ErrEventBadRequest
	}

	return &event, nil
}

func (r *EventPostgresRepository) AllByNameAndSource(ctx context.Context, name *types.EventName, source *types.EventSource) ([]*types.Event, error) {
	var events []*types.Event

	err := r.db.WithContext(ctx).
		Where("name = ? AND source = ?", name.Value(), source.Value()).
		Order("created_at DESC").
		Find(&events).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrEventNotFound
		}
		return nil, shared.ErrEventBadRequest
	}

	return events, nil
}

func (r *EventPostgresRepository) Delete(ctx context.Context, id *types.SharedId) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id.Value()).
		Delete(&types.Event{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return shared.ErrEventNotFound
	}

	return nil
}

func (r *EventPostgresRepository) CreateSubscription(ctx context.Context, sub *types.Subscription) error {
	return r.db.WithContext(ctx).
		Table("subscriptions").
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(sub).Error
}
