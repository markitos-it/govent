package database

import (
	"context"
	"fmt"
	"go-vents/internal/domain/shared"
	"go-vents/internal/domain/types"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
)

type EventSpannerRepository struct {
	client *spanner.Client
}

func NewEventSpannerRepository(client *spanner.Client) types.EventRepository {
	return &EventSpannerRepository{client: client}
}

func (r *EventSpannerRepository) Create(ctx context.Context, event *types.Event) error {
	_, err := r.client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {

		eventStmt := spanner.Statement{
			SQL: `INSERT INTO events (id, name, source, payload, created_at, updated_at) 
			      VALUES (@id, @name, @source, @payload, @created_at, @updated_at)`,
			Params: map[string]interface{}{
				"id":         event.Id,
				"name":       event.Name,
				"source":     event.Source,
				"payload":    event.Payload,
				"created_at": event.CreatedAt,
				"updated_at": event.UpdatedAt,
			},
		}
		if _, err := tx.Update(ctx, eventStmt); err != nil {
			return fmt.Errorf("error inserting main event: %w", err)
		}

		subStmt := spanner.Statement{
			SQL: "SELECT subscriber_name FROM subscriptions WHERE event_name = @name AND source = @source",
			Params: map[string]interface{}{
				"name":   event.Name,
				"source": event.Source,
			},
		}
		iter := tx.Query(ctx, subStmt)
		defer iter.Stop()

		var subscribers []string
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("error searching subscriptions: %w", err)
			}
			var subscriberName string
			if err := row.Columns(&subscriberName); err != nil {
				return err
			}
			subscribers = append(subscribers, subscriberName)
		}

		for _, subscriber := range subscribers {
			queueId := uuid.New().String()
			msg, err := types.NewQueueMessage(queueId, subscriber, event.Id)
			if err != nil {
				return fmt.Errorf("error al construir QueueMessage: %w", err)
			}

			queueStmt := spanner.Statement{
				SQL: `INSERT INTO queue (id, subscriber, event_id, status, created_at, updated_at) 
				      VALUES (@id, @subscriber, @event_id, @status, @created_at, @updated_at)`,
				Params: map[string]interface{}{
					"id":         msg.Id,
					"subscriber": msg.SubscriberName,
					"event_id":   msg.EventId,
					"status":     msg.Status,
					"created_at": msg.CreatedAt,
					"updated_at": msg.UpdatedAt,
				},
			}
			if _, err := tx.Update(ctx, queueStmt); err != nil {
				return fmt.Errorf("error al insertar en la cola de %s: %w", subscriber, err)
			}
		}

		return nil
	})

	return err
}

func (r *EventSpannerRepository) One(ctx context.Context, id *types.SharedId) (*types.Event, error) {
	row, err := r.client.Single().ReadRow(ctx, "events", spanner.Key{id.Value()},
		[]string{"id", "name", "source", "payload", "created_at", "updated_at"})

	if err != nil {
		if spanner.ErrCode(err) == codes.NotFound {
			return nil, shared.ErrEventNotFound
		}
		return nil, shared.ErrEventBadRequest
	}

	var event types.Event
	err = row.Columns(&event.Id, &event.Name, &event.Source, &event.Payload, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventSpannerRepository) AllByNameAndSource(ctx context.Context, name *types.Name, source *types.Source) ([]*types.Event, error) {
	stmt := spanner.Statement{
		SQL:    "SELECT id, name, source, payload, created_at, updated_at FROM events WHERE name = @name AND source = @source ORDER BY created_at DESC",
		Params: map[string]interface{}{"name": name.Value(), "source": source.Value()},
	}

	iter := r.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var events []*types.Event
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, shared.ErrEventBadRequest
		}

		var event types.Event
		err = row.Columns(&event.Id, &event.Name, &event.Source, &event.Payload, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if len(events) == 0 {
		return nil, shared.ErrEventNotFound
	}

	return events, nil
}

func (r *EventSpannerRepository) Delete(ctx context.Context, id *types.SharedId) error {
	_, err := r.client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		_, err := tx.ReadRow(ctx, "events", spanner.Key{id.Value()}, []string{"id"})
		if err != nil {
			if spanner.ErrCode(err) == codes.NotFound {
				return shared.ErrEventNotFound
			}
			return err
		}

		// Si existe, aplicamos la mutación de borrado
		mutation := spanner.Delete("events", spanner.Key{id.Value()})
		return tx.BufferWrite([]*spanner.Mutation{mutation})
	})

	return err
}

func (r *EventSpannerRepository) CreateSubscription(ctx context.Context, sub *types.Subscription) error {
	mutation := spanner.InsertOrUpdate(
		"subscriptions",
		[]string{"id", "subscriber_name", "event_name", "source", "created_at", "updated_at"},
		[]interface{}{sub.Id, sub.SubscriberName, sub.EventName, sub.Source, sub.CreatedAt, sub.UpdatedAt},
	)

	_, err := r.client.Apply(ctx, []*spanner.Mutation{mutation})
	return err
}

func (r *EventSpannerRepository) PullMessages(ctx context.Context, eventName *types.Name, source *types.Source) ([]*types.QueueMessage, error) {
	stmt := spanner.Statement{
		SQL: `SELECT q.id, q.subscriber, q.event_id, q.status, q.created_at, q.updated_at 
		      FROM queue q
		      JOIN events e ON q.event_id = e.id
		      WHERE e.name = @name AND e.source = @source AND q.status = @status`,
		Params: map[string]interface{}{
			"name":   eventName.Value(),
			"source": source.Value(),
			"status": "pending",
		},
	}

	iter := r.client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var results []*types.QueueMessage
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var msg types.QueueMessage
		err = row.Columns(&msg.Id, &msg.SubscriberName, &msg.EventId, &msg.Status, &msg.CreatedAt, &msg.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, &msg)
	}

	return results, nil
}

func (r *EventSpannerRepository) AckMessage(ctx context.Context, id *types.SharedId) error {
	_, err := r.client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		row, err := tx.ReadRow(ctx, "queue", spanner.Key{id.Value()}, []string{"status"})
		if err != nil {
			if spanner.ErrCode(err) == codes.NotFound {
				return shared.ErrQueueMessageNotFound
			}
			return err
		}

		var currentStatus string
		if err := row.Columns(&currentStatus); err != nil {
			return err
		}

		if currentStatus != "pending" {
			return shared.ErrQueueMessageNotFound
		}

		mutation := spanner.Update(
			"queue",
			[]string{"id", "status", "updated_at"},
			[]interface{}{id.Value(), "processed", time.Now()},
		)

		return tx.BufferWrite([]*spanner.Mutation{mutation})
	})

	return err
}
