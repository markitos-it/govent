package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"govent/internal/domain/types"
	"govent/internal/infrastructure/configuration"
	"govent/internal/infrastructure/database"
	"govent/internal/infrastructure/gapi"
	slogcolored "govent/internal/infrastructure/logging/slog-colored"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var repository types.EventRepository
var config configuration.EventConfiguration

// #[.'.]:> Main function that orchestrates the startup and controlled shutdown of the application
// #[.'.]:> STEP 1: Show startup banner
// #[.'.]:> These logs help clearly identify the service startup
// #[.'.]:> STEP 2: Load configuration from files or environment variables
// #[.'.]:> This function sets all operational parameters
// #[.'.]:> STEP 3: Initialize database connection and repository
// #[.'.]:> Prepare data access and table structure
// #[.'.]:> STEP 4: Start gRPC servers
// #[.'.]:> Start entry points for gRPC clients
// #[.'.]:> STEP 5: Show shutdown banner when finished
// #[.'.]:> These logs clearly mark the end of the service execution
func main() {
	logger := slogcolored.NewColoredSLogger()
	logger.OpenGroup("main")

	logger.Info("['.']:>--------------------------------------------")
	logger.Info("['.']:>--- <starting markitos-it-svc-event>  ---")
	logger.Info("['.']:>------- logger loaded")
	loadConfiguration(logger)
	loadDatabase(logger)
	startServers(logger)
	logger.Info("['.']:>--------------------------------------------")
	logger.Info("['.']:>--- <markitos-it-svc-event stopped>  ---")
	logger.Info("['.']:>")
}

// #[.'.]:> This function loads the service configuration
// #[.'.]:> STEP 1: Try to load configuration from file or environment variables
// #[.'.]:> Looks for "app.env" in the current directory, or uses environment variables if not found
// #[.'.]:> If there's an error, terminate the application immediately
// #[.'.]:> Can't operate without valid configuration
// #[.'.]:> STEP 2: Store configuration in a global variable
// #[.'.]:> Makes it accessible to the rest of the program functions
func loadConfiguration(logger types.Logger) {
	loadedConfig, err := configuration.LoadConfiguration(".", logger)
	if err != nil {
		logger.Error("['.']:>------- unable to load configuration: " + err.Error())
		os.Exit(1)
	}

	config = loadedConfig
	logger.Info("['.']:>------- configuration loaded")

}

// #[.'.]:> This function initializes the database and repository
// #[.'.]:> STEP 1: Establish connection to PostgreSQL using the connection string
// #[.'.]:> GORM abstracts connection details and database handling
// #[.'.]:> If unable to connect to the database, it's a fatal error
// #[.'.]:> STEP 2: Run automatic migrations to create or update tables
// #[.'.]:> Ensures the database structure matches our models
// #[.'.]:> If migrations fail, can't continue
// #[.'.]:> STEP 3: Create a repository instance with the database connection
// #[.'.]:> The repository encapsulates all data access logic
func loadDatabase(logger types.Logger) {
	customLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(config.DatabaseDsn), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		logger.Fatal("['.']:> error unable to connect to database:" + err.Error())
	}

	err = db.AutoMigrate(&types.Event{})
	if err != nil {
		logger.Fatal("['.']:> error unable to migrate database:" + err.Error())
	}

	repo := database.NewEventPostgresRepository(db)
	repository = repo

	logger.Info("['.']:>------- database initialized")
}

// #[.'.]:> This function starts the servers and manages their lifecycle
// #[.'.]:> STEP 1: Create a cancelable context to signal shutdown
// #[.'.]:> This context will be propagated to the servers to manage their lifecycle
// #[.'.]:> STEP 2: Set up a channel to capture OS signals
// #[.'.]:> Allows detection of Ctrl+C or system shutdown signals
// #[.'.]:> STEP 3: Create a wait group to coordinate server shutdown
// #[.'.]:> The WaitGroup lets us wait for both servers to fully stop
// #[.'.]:> STEP 4: Start the gRPC server in a separate goroutine
// #[.'.]:> STEP 5: Block until a termination signal is received
// #[.'.]:> The application will wait here until Ctrl+C or SIGTERM is received
// #[.'.]:> STEP 6: Cancel the context to start controlled shutdown
// #[.'.]:> This sends the termination signal to both servers
// #[.'.]:> STEP 7: Wait for both servers to fully stop
// #[.'.]:> Won't exit until both servers have completed their shutdown
func startServers(logger types.Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := runGRPCServer(ctx, logger); err != nil {
			logger.Fatal("['.']:> error running gRPC server: " + err.Error())
		}
	}()
	<-stop
	logger.Info("['.']:>------- shutting down servers gracefully...")
	cancel()
	wg.Wait()
}

// #[.'.]:> This function starts and manages the gRPC server lifecycle
// #[.'.]:> STEP 1: Create a network listener
// #[.'.]:> This listener will listen for TCP requests at the configured address and port
// #[.'.]:> STEP 2: Build a generic unary interceptor to audit payloads (input/output) and execution times
// #[.'.]:> STEP 3: Create a new gRPC server instance injecting the interceptor
// #[.'.]:> STEP 4: Create the implementation of our service
// #[.'.]:> This part contains the actual business logic
// #[.'.]:> STEP 5: Register our service with the gRPC server
// #[.'.]:> Connects our implementations with the gRPC system
// #[.'.]:> STEP 6: Enable reflection to facilitate testing
// #[.'.]:> Reflection allows tools like grpcurl to discover our services
// #[.'.]:> STEP 7: Set up controlled (graceful) shutdown
// #[.'.]:> This goroutine runs in the background and waits for the shutdown signal
// #[.'.]:> Blocks until the context is canceled (shutdown signal)
// #[.'.]:> Logs a message indicating the server is shutting down
// #[.'.]:> Performs a graceful shutdown:
// #[.'.]:> - Stops accepting new connections
// #[.'.]:> - Waits for ongoing calls to finish
// #[.'.]:> - Closes all connections cleanly
// #[.'.]:> STEP 8: Log that the server is running
// #[.'.]:> STEP 9: Start the server (this method blocks until an error occurs)
// #[.'.]:> The server now actively listens for incoming requests
func runGRPCServer(ctx context.Context, logger types.Logger) error {
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		return err
	}

	genericUnaryInterceptor := func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startTime := time.Now()

		var reqJSON string
		if protoReq, ok := req.(proto.Message); ok {
			reqJSON = protojson.MarshalOptions{EmitUnpopulated: true}.Format(protoReq)
		} else {
			reqJSON = fmt.Sprintf("%v", req)
		}

		logger.Info(fmt.Sprintf("gRPC start call ➜ %s | input: %s", info.FullMethod, reqJSON))

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		var respJSON string
		if err == nil {
			if protoResp, ok := resp.(proto.Message); ok {
				respJSON = protojson.MarshalOptions{EmitUnpopulated: true}.Format(protoResp)
			} else {
				respJSON = fmt.Sprintf("%v", resp)
			}
		}

		if err != nil {
			logger.Error(fmt.Sprintf("gRPC finish call ❌ %s | duration: %v | error: %v", info.FullMethod, duration, err))
		} else {
			logger.Info(fmt.Sprintf("gRPC finish call  %s | duration: %v | output: %s", info.FullMethod, duration, respJSON))
		}

		return resp, err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(genericUnaryInterceptor),
	)

	server := gapi.NewServer(config.GRPCServerAddress, repository, config, logger)
	gapi.RegisterEventserviceServer(grpcServer, server)
	reflection.Register(grpcServer)

	go func() {
		<-ctx.Done()
		logger.Info("['.']:> shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	logger.Info("['.']:> gRPC server running at address: " + config.GRPCServerAddress)
	logger.CloseGroup("main")

	return grpcServer.Serve(listener)
}
