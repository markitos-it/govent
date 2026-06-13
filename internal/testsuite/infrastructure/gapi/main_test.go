package gapi_test

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"go-vents/internal/domain/types"
	"go-vents/internal/infrastructure/configuration"
	"go-vents/internal/infrastructure/gapi"
	slogcolored "go-vents/internal/infrastructure/logging/slog-colored"
	"go-vents/internal/testsuite/infrastructure/testdb"
	internal_test "go-vents/internal/testsuite/internal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var grpcServer *grpc.Server
var grpcClient gapi.EventserviceClient
var ctx context.Context

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	grpcServer.Stop()
	os.Exit(code)
}

func setup() {
	lis = bufconn.Listen(bufSize)

	grpcServer = grpc.NewServer()

	config := &configuration.EventConfiguration{}
	server := gapi.NewServer(":8080", testdb.GetRepository(), *config, slogcolored.NewColoredSLogger())
	gapi.RegisterEventserviceServer(grpcServer, server)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("['.']:> Error serving gRPC server: %v", err)
		}
	}()

	conn, err := grpc.NewClient(
		"passthrough://localhost/bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("['.']:> Failed to dial bufnet: %v", err)
	}

	grpcClient = gapi.NewEventserviceClient(conn)
	ctx = context.Background()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func createPersistedRandomEvent() *types.Event {
	event := internal_test.NewRandomEvent()
	_ = testdb.GetRepository().Create(ctx, event)

	return event
}

func deletePersistedRandomEvent(eventId string) {
	id, err := types.NewSharedId(eventId)
	if err != nil {
		return
	}

	_ = testdb.GetRepository().Delete(ctx, id)
}
