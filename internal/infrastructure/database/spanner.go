package database

import (
	"context"
	"errors"
	"go-vents/internal/domain/shared"
	"go-vents/internal/domain/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventSpannerRepository struct {
	db *gorm.DB
}

func NewEventSpannerRepository(db *gorm.DB) types.EventRepository {
	return &EventSpannerRepository{db: db}
}

func (r *EventSpannerRepository) Create(ctx context.Context, event *types.Event) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		var subscribers []string
		err := tx.Table("subscriptions").
			Where("event_name = ? AND source = ?", event.Slug, event.Source).
			Pluck("subscriber_name", &subscribers).Error
		if err != nil {
			return err
		}

		for _, subscriber := range subscribers {
			queueId := uuid.New().String()
			msg, err := types.NewQueueMessage(queueId, subscriber, event.Id)
			if err != nil {
				return err
			}
			if err := tx.Table("queue").Create(msg).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *EventSpannerRepository) One(ctx context.Context, id *types.SharedId) (*types.Event, error) {
	var event types.Event
	err := r.db.WithContext(ctx).First(&event, "id = ?", id.Value()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrEventNotFound
		}
		return nil, shared.ErrEventBadRequest
	}
	return &event, nil
}

func (r *EventSpannerRepository) AllBySlugAndSource(ctx context.Context, slug *types.Slug, source *types.Source) ([]*types.Event, error) {
	var events []*types.Event
	err := r.db.WithContext(ctx).
		Where("slug = ? AND source = ?", slug.Value(), source.Value()).
		Order("created_at DESC").
		Find(&events).Error
	return events, err
}

func (r *EventSpannerRepository) Delete(ctx context.Context, id *types.SharedId) error {
	result := r.db.WithContext(ctx).Where("id = ?", id.Value()).Delete(&types.Event{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return shared.ErrEventNotFound
	}
	return nil
}

func (r *EventSpannerRepository) CreateSubscription(ctx context.Context, sub *types.Subscription) error {
	return r.db.WithContext(ctx).
		Table("subscriptions").
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(sub).Error
}

func (r *EventSpannerRepository) PullMessages(ctx context.Context, slug *types.Slug, source *types.Source) ([]*types.Queue, error) {
	var results []*types.Queue
	err := r.db.WithContext(ctx).
		Table("queue q").
		Joins("JOIN events e ON q.event_id = e.id").
		Where("e.slug = ? AND e.source = ? AND q.status = ?", slug.Value(), source.Value(), "pending").
		Find(&results).Error
	return results, err
}

func (r *EventSpannerRepository) AckMessage(ctx context.Context, id *types.SharedId) error {
	result := r.db.WithContext(ctx).
		Model(&types.Queue{}).
		Where("id = ? AND status = ?", id.Value(), "pending").
		Updates(map[string]interface{}{
			"status":     "processed",
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return shared.ErrQueueMessageNotFound
	}
	return nil
}
