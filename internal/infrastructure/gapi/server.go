package gapi

import (
	"errors"
	"strings"

	sharedDomain "go-vents/internal/domain/shared"
	"go-vents/internal/domain/types"
	"go-vents/internal/infrastructure/configuration"

	codes "google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	UnimplementedEventserviceServer
	address    string
	repository types.EventRepository
	config     configuration.EventConfiguration
	logger     types.Logger
}

func NewServer(address string, repository types.EventRepository, config configuration.EventConfiguration, logger types.Logger) *Server {
	apiGRPC := &Server{
		address:    address,
		repository: repository,
		config:     config,
		logger:     logger,
	}

	return apiGRPC
}

func (s *Server) GetGRPCCode(err error) codes.Code {
	var code = codes.Internal

	switch {
	case errors.Is(err, sharedDomain.ErrEventNotFound):
		code = codes.NotFound
	case errors.Is(err, sharedDomain.ErrInvalidEventId),
		errors.Is(err, sharedDomain.ErrInvalidEventName),
		errors.Is(err, sharedDomain.ErrInvalidPageNumber),
		errors.Is(err, sharedDomain.ErrInvalidPageSize),
		strings.Contains(err.Error(), "invalid"),
		strings.Contains(err.Error(), "Invalid"),
		strings.Contains(err.Error(), "illegal"),
		strings.Contains(err.Error(), "bad request"):
		code = codes.InvalidArgument
	}

	return code
}

func (s *Server) ToProtos(domainEvents []*types.Event) []*Event {
	var protoEvents []*Event
	for _, event := range domainEvents {
		protoEvents = append(protoEvents, s.ToProto(event))
	}

	return protoEvents
}

func (s *Server) ToProto(domainEvent *types.Event) *Event {
	return &Event{
		Id:        domainEvent.Id,
		Name:      domainEvent.Name,
		Source:    domainEvent.Source,
		Payload:   domainEvent.Payload,
		CreatedAt: timestamppb.New(domainEvent.CreatedAt),
		UpdatedAt: timestamppb.New(domainEvent.UpdatedAt),
	}
}
