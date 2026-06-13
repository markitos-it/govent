package types

import (
	"regexp"

	"go-vents/internal/domain/shared"
)

type Slug struct {
	value string
}

func NewSlug(value string) (*Slug, error) {

	if isValidSlug(value) {
		return &Slug{value}, nil
	}

	return nil, shared.ErrEventBadRequest
}

func isValidSlug(value string) bool {
	if len(value) == 0 {
		return false
	}

	pattern := `^[a-zA-Z]([a-zA-Z0-9._-]*[a-zA-Z0-9])?$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

func (b *Slug) Value() string {
	return b.value
}
