package types

import "context"

type EventRepository interface {
	Create(ctx context.Context, event *Event) error
	One(ctx context.Context, id *SharedId) (*Event, error)
	AllByNameAndSource(ctx context.Context, name *EventName, source *EventSource) ([]*Event, error)
	Delete(ctx context.Context, id *SharedId) error
	CreateSubscription(ctx context.Context, sub *Subscription) error
}
