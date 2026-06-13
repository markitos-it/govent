package types

import (
	"regexp"

	"go-vents/internal/domain/shared"
)

type Name struct {
	value string
}

const EVENT_NAME_MAX_LENGTH = 100
const EVENT_NAME_MIN_LENGTH = 3

func NewName(value string) (*Name, error) {
	if isValidName(value) {
		return &Name{value}, nil
	}

	return nil, shared.ErrEventBadRequest
}

func isValidName(value string) bool {
	if len(value) > EVENT_NAME_MAX_LENGTH || len(value) < EVENT_NAME_MIN_LENGTH {
		return false
	}

	pattern := `^[a-zA-Z0-9]([a-zA-Z0-9 ]*[a-zA-Z0-9])?$`
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false
	}

	return matched
}

func (b *Name) Value() string {
	return b.value
}
