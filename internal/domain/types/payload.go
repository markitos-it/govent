package types

import (
	"encoding/json"

	"go-vents/internal/domain/shared"
)

type Payload struct {
	value string
}

func NewPayload(value string) (*Payload, error) {

	if isValidPayload(value) {
		return &Payload{value}, nil
	}

	return nil, shared.ErrEventBadRequest
}

func isValidPayload(value string) bool {
	if value == "" {
		return true
	}
	return json.Valid([]byte(value))
}

func (b *Payload) Value() string {
	return b.value
}
