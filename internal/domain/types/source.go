package types

import (
	"go-vents/internal/domain/shared"
)

type Source struct {
	value string
}

const EVENT_CONTENT_MAX_LENGTH = 200
const EVENT_CONTENT_MIN_LENGTH = 1

func NewSource(value string) (*Source, error) {

	if isValidSource(value) {
		return &Source{value}, nil
	}

	return nil, shared.ErrEventBadRequest
}

func isValidSource(value string) bool {
	return len(value) >= EVENT_CONTENT_MIN_LENGTH || len(value) <= EVENT_CONTENT_MAX_LENGTH
}

func (b *Source) Value() string {
	return b.value
}
