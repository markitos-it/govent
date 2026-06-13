package types

import (
	"regexp"

	"go-vents/internal/domain/shared"
)

type EventName struct {
	value string
}

const EVENT_NAME_MAX_LENGTH = 100
const EVENT_NAME_MIN_LENGTH = 3

func NewEventName(value string) (*EventName, error) {
	if isValidEventName(value) {
		return &EventName{value}, nil
	}

	return nil, shared.ErrEventBadRequest
}

func isValidEventName(value string) bool {
	if len(value) > EVENT_NAME_MAX_LENGTH || len(value) < EVENT_NAME_MIN_LENGTH {
		return false
	}

	pattern := `^[a-zA-Z]{1}[a-zA-Z ]+[a-zA-Z]$|^[a-zA-Z]{1}$`
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false
	}

	return matched
}

func (b *EventName) Value() string {
	return b.value
}
