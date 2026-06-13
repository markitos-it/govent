package types

import (
	"context"
)

type EventRepository interface {
	Create(ctx context.Context, event *Event) error
	One(ctx context.Context, id *SharedId) (*Event, error)
	AllByNameAndSource(ctx context.Context, name *Name, source *Source) ([]*Event, error)
	Delete(ctx context.Context, id *SharedId) error
	CreateSubscription(ctx context.Context, sub *Subscription) error
	PullMessages(ctx context.Context, name *Name, source *Source) ([]*QueueMessage, error)
	AckMessage(ctx context.Context, id *SharedId) error
}
