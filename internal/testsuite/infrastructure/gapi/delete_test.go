package gapi_test

import (
	"log"
	"testing"

	"go-vents/internal/domain/shared"
	"go-vents/internal/infrastructure/gapi"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEventCanDelete(t *testing.T) {
	event := createPersistedRandomEvent()

	resp, err := grpcClient.DeleteEvent(ctx, &gapi.DeleteEventRequest{
		Id: event.Id,
	})

	require.NoError(t, err)
	require.Equal(t, event.Id, resp.Deleted)

	_, err = grpcClient.GetEvent(ctx, &gapi.GetEventRequest{Id: event.Id})
	require.Error(t, err)
}

func TestEventCantDeleteValidButNonExistingId(t *testing.T) {
	request := &gapi.DeleteEventRequest{
		Id: shared.UUIDv4(),
	}

	ob, err := grpcClient.DeleteEvent(ctx, request)

	log.Println("-----------------------------------------------------")
	log.Println("request", request)
	log.Println("err", err)
	log.Println("ob", ob)
	log.Println("-----------------------------------------------------")

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
}

func TestEventCantDeleteInvalidEventId(t *testing.T) {
	_, err := grpcClient.DeleteEvent(ctx, &gapi.DeleteEventRequest{
		Id: "an-invalid-id-format",
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}
