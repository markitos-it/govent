package gapi_test

import (
	"testing"

	"go-vents/internal/infrastructure/gapi"
	internal_test "go-vents/internal/testsuite/internal"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEventCanCreate(t *testing.T) {
	event := internal_test.NewRandomOnlyNameEvent()

	resp, err := grpcClient.CreateEvent(ctx, &gapi.CreateEventRequest{
		Name: event.Name,
		/* ___CUSTOM_TEST_FIELDS___*/
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp.Id)
	require.Equal(t, event.Name, resp.Name)

	deletePersistedRandomEvent(resp.Id)
}

func TestEventCantCreateWithoutName(t *testing.T) {
	_, err := grpcClient.CreateEvent(ctx, &gapi.CreateEventRequest{
		Name: "",
		/* ___CUSTOM_REQUIRED_VALUES___*/
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}

func TestEventCantCreateWithoutValidName(t *testing.T) {
	_, err := grpcClient.CreateEvent(ctx, &gapi.CreateEventRequest{
		Name: "!!!!!invalid!!!name!!!",
		/* ___CUSTOM_REQUIRED_VALUES___*/
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}
