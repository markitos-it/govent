package services_test

import (
	"context"
	"os"
	"testing"

	"govent/internal/domain/types"
	internal_test "govent/internal/testsuite/internal"
)

type MockSpyEventRepository struct {
	LastCreatedEventName         *string
	LastDeleteEventId            *string
	LastOneEventId               *string
	LastAllByNameAndSource       []types.Event
	LastCreatedSubscriptionEvent *string
}

func NewMockSpyEventRepository() *MockSpyEventRepository {
	return &MockSpyEventRepository{
		LastCreatedEventName:         nil,
		LastDeleteEventId:            nil,
		LastOneEventId:               nil,
		LastAllByNameAndSource:       nil,
		LastCreatedSubscriptionEvent: nil,
	}
}

func (m *MockSpyEventRepository) Create(ctx context.Context, event *types.Event) error {
	m.LastCreatedEventName = &event.Name

	return nil
}

func (m *MockSpyEventRepository) CreateHaveBeenCalledWith(eventName *string) bool {
	var result = m.LastCreatedEventName != nil && *m.LastCreatedEventName == *eventName

	m.LastCreatedEventName = nil

	return result
}

func (m *MockSpyEventRepository) CreateSubscription(ctx context.Context, sub *types.Subscription) error {
	m.LastCreatedSubscriptionEvent = &sub.EventName
	return nil
}

func (m *MockSpyEventRepository) CreateSubscriptionHaveBeenCalledWith(eventName *string) bool {
	var result = m.LastCreatedSubscriptionEvent != nil && *m.LastCreatedSubscriptionEvent == *eventName

	m.LastCreatedSubscriptionEvent = nil

	return result
}

func (m *MockSpyEventRepository) Delete(ctx context.Context, id *types.SharedId) error {
	value := id.Value()
	m.LastDeleteEventId = &value

	return nil
}

func (m *MockSpyEventRepository) DeleteHaveBeenCalledWith(eventId *string) bool {
	var result = m.LastDeleteEventId != nil && *m.LastDeleteEventId == *eventId

	m.LastDeleteEventId = nil

	return result
}

func (m *MockSpyEventRepository) One(ctx context.Context, id *types.SharedId) (*types.Event, error) {
	value := id.Value()
	m.LastOneEventId = &value

	return internal_test.NewRandomEventWithCustomId(id), nil
}

func (m *MockSpyEventRepository) OneHaveBeenCalledWith(eventId *string) bool {
	var result = m.LastOneEventId != nil && *m.LastOneEventId == *eventId

	m.LastOneEventId = nil

	return result
}

func (m *MockSpyEventRepository) AllByNameAndSource(ctx context.Context, name *types.EventName, source *types.EventSource) ([]*types.Event, error) {
	anEvent := internal_test.NewRandomEventWithNameAndSource(name.Value(), source.Value())
	m.LastAllByNameAndSource = append(m.LastAllByNameAndSource, *anEvent)

	return []*types.Event{anEvent}, nil
}

func (m *MockSpyEventRepository) LastAllByNameAndSourceHaveBeenCalled(name *types.EventName, source *types.EventSource) bool {
	var result = m.LastAllByNameAndSource[0].Name == name.Value() &&
		m.LastAllByNameAndSource[0].Source == source.Value()

	m.LastAllByNameAndSource = nil

	return result
}

var repository *MockSpyEventRepository

func TestMain(m *testing.M) {
	repository = NewMockSpyEventRepository()

	os.Exit(m.Run())
}
