package types

import "go-vents/internal/domain/shared"

type SharedId struct {
	value string
}

func NewSharedId(value string) (*SharedId, error) {
	if shared.IsUUIDv4(value) {
		return &SharedId{value}, nil
	}

	return nil, shared.ErrEventBadRequest
}

func (b *SharedId) Value() string {
	return b.value
}
