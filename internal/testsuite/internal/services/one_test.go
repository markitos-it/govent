package services_test

import (
	"context"
	"testing"

	"go-vents/internal/domain/services"
	"go-vents/internal/domain/shared"

	"github.com/stretchr/testify/assert"
)

func TestCanGetAResource(t *testing.T) {
	var request = services.EventOneRequest{
		Id: shared.UUIDv4(),
	}

	var service = services.NewEventOneService(repository)
	event, err := service.Do(context.TODO(), request)

	assert.Nil(t, err)
	assert.IsType(t, services.EventOneResponse{}, *event)
	assert.True(t, repository.OneHaveBeenCalledWith(&request.Id))
}

func TestCantGetWithoutId(t *testing.T) {
	var request = services.EventOneRequest{}

	var service = services.NewEventOneService(repository)
	_, err := service.Do(context.TODO(), request)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrEventBadRequest)
	assert.False(t, repository.OneHaveBeenCalledWith(&request.Id))
}

func TestCantGetWithoutInvalidId(t *testing.T) {
	var request = services.EventOneRequest{
		Id: "invalid-id",
	}

	var service = services.NewEventOneService(repository)
	_, err := service.Do(context.TODO(), request)

	assert.ErrorIs(t, err, shared.ErrEventBadRequest)
	assert.False(t, repository.OneHaveBeenCalledWith(&request.Id))
}
