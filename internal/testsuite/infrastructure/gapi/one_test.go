package gapi_test

import (
	"testing"

	"go-vents/internal/domain/shared"
	"go-vents/internal/infrastructure/gapi"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEventCanGetOne(t *testing.T) {
	event := createPersistedRandomEvent()

	resp, err := grpcClient.GetEvent(ctx, &gapi.GetEventRequest{
		Id: event.Id,
	})

	require.NoError(t, err)
	require.Equal(t, event.Name, resp.Name)
	require.Equal(t, event.Id, resp.Id)

	deletePersistedRandomEvent(resp.Id)
}

func TestEventCantGetInvalidId(t *testing.T) {
	_, err := grpcClient.GetEvent(ctx, &gapi.GetEventRequest{
		Id: "an-invalid-id",
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}

func TestEventCantGetValidIdButNonExistingResource(t *testing.T) {
	_, err := grpcClient.GetEvent(ctx, &gapi.GetEventRequest{
		Id: shared.UUIDv4(),
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
}
