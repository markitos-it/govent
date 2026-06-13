package types

import (
	"encoding/json"

	"go-vents/internal/domain/shared"
)

type EventPayload struct {
	value string
}

func NewEventPayload(value string) (*EventPayload, error) {

	if isValidEventPayload(value) {
		return &EventPayload{value}, nil
	}

	return nil, shared.ErrEventBadRequest
}

func isValidEventPayload(value string) bool {
	if value == "" {
		return true
	}
	return json.Valid([]byte(value))
}

func (b *EventPayload) Value() string {
	return b.value
}
