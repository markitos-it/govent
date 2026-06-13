package domain_test

import (
	"strings"
	"testing"

	"go-vents/internal/domain/shared"
	"go-vents/internal/domain/types"

	"github.com/stretchr/testify/assert"
)

func TestCanCreateValidEventName(t *testing.T) {
	validNames := []string{
		"ValidName",
		"AnotherValidName",
		"Valid Name With Spaces",
		"Short",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"}
	for _, name := range validNames {
		if _, err := types.NewEventName(name); err != nil {
			t.Errorf("Expected valid name, but got invalid: %s", name)
		}
	}

	invalidNames := []string{
		" InvalidName",
		"InvalidName ",
		"Invalid Name ",
		" Invalid Name",
		"Invalid@Name",
		"Invalid#Name",
		"Invalid123Name",
		"Invalid Name With Spaces ",
		" Invalid Name With Spaces",
		"Invalid Name With Spaces And Symbols!",
	}
	for _, name := range invalidNames {
		if _, err := types.NewEventName(name); err == nil {
			t.Errorf("Expected valid name, but got invalid: %s", name)
		}
	}

	invalidLengthNames := []string{
		strings.Repeat("a", types.EVENT_NAME_MAX_LENGTH+1),
		strings.Repeat("b", types.EVENT_NAME_MIN_LENGTH-1),
		"",
	}
	for _, name := range invalidLengthNames {
		if _, err := types.NewEventName(name); err == nil {
			t.Errorf("Expected invalid name, but got invalid: %s", name)
		}
	}
}

func TestCanCreateValidEventPayload(t *testing.T) {
	validPayloads := []string{
		"",
		"{}",
		`{"key": "value"}`,
		`[1, 2, 3]`,
		`"just a string"`,
		"123",
	}

	for _, payload := range validPayloads {
		if _, err := types.NewEventPayload(payload); err != nil {
			t.Errorf("Expected valid payload, but got error for: %s", payload)
		}
	}

	invalidPayloads := []string{
		`{bad json}`,
		`{"key": "value",}`,
		`[1, 2,, 3]`,
		`"unclosed string`,
	}

	for _, payload := range invalidPayloads {
		if _, err := types.NewEventPayload(payload); err == nil {
			t.Errorf("Expected invalid payload, but got valid for: %s", payload)
		}
	}
}

func TestCanCreateValidQueueMessage(t *testing.T) {
	id := shared.UUIDv4()
	subscriberName := "sub1"
	eventId := shared.UUIDv4()

	qm, err := types.NewQueueMessage(id, subscriberName, eventId)
	assert.Nil(t, err)
	assert.NotNil(t, qm)
	assert.Equal(t, id, qm.Id)
	assert.Equal(t, subscriberName, qm.SubscriberName)
	assert.Equal(t, eventId, qm.EventId)
	assert.Equal(t, types.StatusPending, qm.Status)
	assert.NotZero(t, qm.CreatedAt)
	assert.NotZero(t, qm.UpdatedAt)
	assert.Equal(t, "queue", qm.TableName())
}

func TestCantCreateQueueMessageWithEmptyFields(t *testing.T) {
	id := shared.UUIDv4()
	subscriberName := "sub1"
	eventId := shared.UUIDv4()

	_, err := types.NewQueueMessage("", subscriberName, eventId)
	assert.NotNil(t, err)

	_, err = types.NewQueueMessage(id, "", eventId)
	assert.NotNil(t, err)

	_, err = types.NewQueueMessage(id, subscriberName, "")
	assert.NotNil(t, err)
}

func TestQueueMessageMarkAsProcessedAndFailed(t *testing.T) {
	id := shared.UUIDv4()
	subscriberName := "sub1"
	eventId := shared.UUIDv4()

	qm, _ := types.NewQueueMessage(id, subscriberName, eventId)

	qm.MarkAsProcessed()
	assert.Equal(t, types.StatusProcessed, qm.Status)

	qm.MarkAsFailed()
	assert.Equal(t, types.StatusFailed, qm.Status)
}

func TestCanCreateValidSubscription(t *testing.T) {
	id := shared.UUIDv4()
	subscriberName := "sub1"
	eventName := "event1"
	source := "source1"

	sub, err := types.NewSubscription(id, subscriberName, eventName, source)
	assert.Nil(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, id, sub.Id)
	assert.Equal(t, subscriberName, sub.SubscriberName)
	assert.Equal(t, eventName, sub.EventName)
	assert.Equal(t, source, sub.Source)
	assert.Equal(t, "subscriptions", sub.TableName())
}

func TestCantCreateSubscriptionWithInvalidFields(t *testing.T) {
	id := shared.UUIDv4()
	subscriberName := "sub1"
	eventName := "event1"
	source := "source1"

	_, err := types.NewSubscription("invalid-uuid", subscriberName, eventName, source)
	assert.NotNil(t, err)

	_, err = types.NewSubscription(id, "", eventName, source)
	assert.NotNil(t, err)

	_, err = types.NewSubscription(id, subscriberName, "", source)
	assert.NotNil(t, err)

	_, err = types.NewSubscription(id, subscriberName, eventName, "")
	assert.NotNil(t, err)
}
