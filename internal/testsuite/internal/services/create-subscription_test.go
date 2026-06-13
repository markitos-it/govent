package services_test

import (
	"context"
	"testing"

	"go-vents/internal/domain/services"

	"github.com/stretchr/testify/assert"
)

func TestCanCreateASubscription(t *testing.T) {
	var request = services.SubscriptionCreateRequest{
		SubscriberName: "Sub1",
		Source:         "Source1",
		Event:          "Event1",
	}

	var service = services.NewSubscriptionCreateService(repository)
	response, err := service.Do(context.TODO(), request)

	assert.Nil(t, err)
	assert.True(t, repository.CreateSubscriptionHaveBeenCalledWith(&request.Event))
	assert.Equal(t, response.Event, request.Event)
	assert.Equal(t, response.SubscriberName, request.SubscriberName)
	assert.Equal(t, response.Source, request.Source)
	assert.NotEmpty(t, response.Id)
}

func TestCantCreateSubscriptionWithoutSubscriberName(t *testing.T) {
	var request = services.SubscriptionCreateRequest{
		Source: "Source1",
		Event:  "Event1",
	}

	var service = services.NewSubscriptionCreateService(repository)
	_, err := service.Do(context.TODO(), request)

	assert.NotNil(t, err)
	assert.False(t, repository.CreateSubscriptionHaveBeenCalledWith(&request.Event))
}

func TestCantCreateSubscriptionWithoutEvent(t *testing.T) {
	var request = services.SubscriptionCreateRequest{
		SubscriberName: "Sub1",
		Source:         "Source1",
	}

	var service = services.NewSubscriptionCreateService(repository)
	_, err := service.Do(context.TODO(), request)

	assert.NotNil(t, err)
	assert.False(t, repository.CreateSubscriptionHaveBeenCalledWith(&request.Event))
}

func TestCantCreateSubscriptionWithoutSource(t *testing.T) {
	var request = services.SubscriptionCreateRequest{
		SubscriberName: "Sub1",
		Event:          "Event1",
	}

	var service = services.NewSubscriptionCreateService(repository)
	_, err := service.Do(context.TODO(), request)

	assert.NotNil(t, err)
	assert.False(t, repository.CreateSubscriptionHaveBeenCalledWith(&request.Event))
}
