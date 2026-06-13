package database

import (
	"context"
	"encoding/json"
	"fmt"
	"go-vents/internal/domain/types"
	"time"

	"github.com/redis/go-redis/v9"
)

type EventCachedRepository struct {
	realRepo types.EventRepository // Spanner, MariaDB o Postgres
	rdb      *redis.Client
	ttl      time.Duration
}

func NewEventCachedRepository(realRepo types.EventRepository, rdb *redis.Client, ttl time.Duration) types.EventRepository {
	return &EventCachedRepository{
		realRepo: realRepo,
		rdb:      rdb,
		ttl:      ttl,
	}
}

func (r *EventCachedRepository) Create(ctx context.Context, event *types.Event) error {
	if err := r.realRepo.Create(ctx, event); err != nil {
		return err
	}

	go func() {
		bgCtx := context.Background()
		cacheKeyCollection := fmt.Sprintf("events:collection:%s:%s", event.Slug, event.Source)
		_ = r.rdb.Del(bgCtx, cacheKeyCollection)

		cacheKeyIndividual := fmt.Sprintf("event:%s", event.Id)
		if jsonData, err := json.Marshal(event); err == nil {
			_ = r.rdb.Set(bgCtx, cacheKeyIndividual, jsonData, r.ttl)
		}
	}()

	return nil
}

func (r *EventCachedRepository) One(ctx context.Context, id *types.SharedId) (*types.Event, error) {
	cacheKey := fmt.Sprintf("event:%s", id.Value())

	val, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var event types.Event
		if err := json.Unmarshal([]byte(val), &event); err == nil {
			return &event, nil
		}
	}

	event, err := r.realRepo.One(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		if jsonData, err := json.Marshal(event); err == nil {
			_ = r.rdb.Set(context.Background(), cacheKey, jsonData, r.ttl).Err()
		}
	}()

	return event, nil
}

func (r *EventCachedRepository) AllBySlugAndSource(ctx context.Context, slug *types.Slug, source *types.Source) ([]*types.Event, error) {
	cacheKey := fmt.Sprintf("events:collection:%s:%s", slug.Value(), source.Value())

	val, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var events []*types.Event
		if err := json.Unmarshal([]byte(val), &events); err == nil {
			return events, nil
		}
	}

	events, err := r.realRepo.AllBySlugAndSource(ctx, slug, source)
	if err != nil {
		return nil, err
	}

	go func() {
		if jsonData, err := json.Marshal(events); err == nil {
			_ = r.rdb.Set(context.Background(), cacheKey, jsonData, r.ttl).Err()
		}
	}()

	return events, nil
}

func (r *EventCachedRepository) Delete(ctx context.Context, id *types.SharedId) error {
	event, err := r.realRepo.One(ctx, id)
	if err != nil {
		return err
	}

	if err := r.realRepo.Delete(ctx, id); err != nil {
		return err
	}

	go func() {
		bgCtx := context.Background()
		cacheKeyIndividual := fmt.Sprintf("event:%s", id.Value())
		cacheKeyCollection := fmt.Sprintf("events:collection:%s:%s", event.Slug, event.Source)

		_ = r.rdb.Del(bgCtx, cacheKeyIndividual, cacheKeyCollection)
	}()

	return nil
}

func (r *EventCachedRepository) CreateSubscription(ctx context.Context, sub *types.Subscription) error {
	return r.realRepo.CreateSubscription(ctx, sub)
}

func (r *EventCachedRepository) PullMessages(ctx context.Context, slug *types.Slug, source *types.Source) ([]*types.Queue, error) {
	return r.realRepo.PullMessages(ctx, slug, source)
}

func (r *EventCachedRepository) AckMessage(ctx context.Context, id *types.SharedId) error {
	return r.realRepo.AckMessage(ctx, id)
}
