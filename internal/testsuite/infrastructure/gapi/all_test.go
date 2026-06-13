package gapi_test

import (
	"testing"

	"go-vents/internal/infrastructure/gapi"

	"github.com/stretchr/testify/require"
)

func TestEventCanListAllResources(t *testing.T) {
	event1 := createPersistedRandomEvent()
	event2 := createPersistedRandomEvent()

	resp, err := grpcClient.AllByNameAndSource(ctx, &gapi.AllEventsByNameAndSourceRequest{
		Name:   event1.Name,
		Source: event1.Source,
	})

	require.NoError(t, err)
	require.NotNil(t, resp.Events)
	found1, found2 := false, false
	for _, b := range resp.Events {
		if b.Id == event1.Id {
			found1 = true
		}
		if b.Id == event2.Id {
			found2 = true
		}
	}
	require.True(t, found1, "First event not found in response")
	require.False(t, found2, "Second event not found in response")

	resp, err = grpcClient.AllByNameAndSource(ctx, &gapi.AllEventsByNameAndSourceRequest{
		Name:   event2.Name,
		Source: event2.Source,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Events)
	found1, found2 = false, false
	for _, b := range resp.Events {
		if b.Id == event1.Id {
			found1 = true
		}
		if b.Id == event2.Id {
			found2 = true
		}
	}
	require.False(t, found1, "First event not found in response")
	require.True(t, found2, "Second event not found in response")

	deletePersistedRandomEvent(event1.Id)
	deletePersistedRandomEvent(event2.Id)
}
