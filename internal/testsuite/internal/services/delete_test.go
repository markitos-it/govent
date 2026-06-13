package services_test

import (
	"context"
	"testing"

	"go-vents/internal/domain/services"
	"go-vents/internal/domain/shared"

	"github.com/stretchr/testify/assert"
)

func TestCanDeleteAUser(t *testing.T) {
	var request = services.EventDeleteRequest{
		Id: shared.UUIDv4(),
	}

	var service = services.NewEventDeleteService(repository)
	err := service.Do(context.TODO(), request)
	assert.Nil(t, err)
	assert.True(t, repository.DeleteHaveBeenCalledWith(&request.Id))
}

func TestCantDeleteWithoutId(t *testing.T) {
	var request = services.EventDeleteRequest{}

	var service = services.NewEventDeleteService(repository)
	err := service.Do(context.TODO(), request)

	assert.ErrorIs(t, err, shared.ErrEventBadRequest)
	assert.NotNil(t, err)
	assert.False(t, repository.DeleteHaveBeenCalledWith(&request.Id))
}

func TestCantDeleteWithInvalidId(t *testing.T) {
	var request = services.EventDeleteRequest{
		Id: "invalid-id",
	}

	var service = services.NewEventDeleteService(repository)
	err := service.Do(context.TODO(), request)

	assert.ErrorIs(t, err, shared.ErrEventBadRequest)
	assert.False(t, repository.DeleteHaveBeenCalledWith(&request.Id))
}
