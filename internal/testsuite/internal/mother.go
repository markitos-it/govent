package internal_test

import (
	"go-vents/internal/domain/shared"
	"go-vents/internal/domain/types"
)

func NewRandomEvent() *types.Event {
	event, _ := types.NewEvent(
		shared.UUIDv4(),
		shared.RandomString(),
		shared.RandomString(),
		"",
	)

	return event
}

func NewRandomOnlyNameEvent() *types.Event {
	event, _ := types.NewEvent(shared.UUIDv4(), shared.RandomString(), "", "")

	return event
}
func NewRandomEventWithNameAndSource(name, source string) *types.Event {
	event, _ := types.NewEvent(
		shared.UUIDv4(),
		name,
		source,
		"",
	)

	return event
}

func NewRandomEventWithCustomId(eventId *types.SharedId) *types.Event {
	event, _ := types.NewEvent(
		eventId.Value(),
		shared.RandomString(),
		shared.RandomString(),
		"",
	)
	return event
}
