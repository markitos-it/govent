package types

import (
	"go-vents/internal/domain/shared"
)

type EventSource struct {
	value string
}

const EVENT_CONTENT_MAX_LENGTH = 200
const EVENT_CONTENT_MIN_LENGTH = 1

func NewEventSource(value string) (*EventSource, error) {

	if isValidEventSource(value) {
		return &EventSource{value}, nil
	}

	return nil, shared.ErrEventBadRequest
}

func isValidEventSource(value string) bool {
	return len(value) >= EVENT_CONTENT_MIN_LENGTH || len(value) <= EVENT_CONTENT_MAX_LENGTH
}

func (b *EventSource) Value() string {
	return b.value
}
