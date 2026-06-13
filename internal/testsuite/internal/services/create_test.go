package services_test

import (
	"context"
	"testing"

	"go-vents/internal/domain/services"
	"go-vents/internal/domain/shared"
	"go-vents/internal/domain/types"

	"github.com/stretchr/testify/assert"
)

func TestCanCreateAUser(t *testing.T) {
	var event = types.Event{
		Name: shared.RandomSlug(),
	}
	var request = services.EventCreateRequest{
		Name: event.Name,
	}

	var service = services.NewEventCreateService(repository)
	response, err := service.Do(context.TODO(), request)

	assert.Nil(t, err)
	assert.True(t, repository.CreateHaveBeenCalledWith(&request.Name))
	assert.Equal(t, response.Name, request.Name)
	assert.NotEmpty(t, response.Id)
}

func TestCantCreateWithoutName(t *testing.T) {
	var request = services.EventCreateRequest{}

	var service = services.NewEventCreateService(repository)
	_, err := service.Do(context.TODO(), request)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrEventBadRequest)
	assert.False(t, repository.CreateHaveBeenCalledWith(&request.Name))
}

func TestCantCreateWithoutValidName(t *testing.T) {
	var request = services.EventCreateRequest{
		Name: "",
	}

	var service = services.NewEventCreateService(repository)
	_, err := service.Do(context.TODO(), request)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrEventBadRequest)
	assert.False(t, repository.CreateHaveBeenCalledWith(&request.Name))
}
