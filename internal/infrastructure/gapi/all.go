package gapi

import (
	context "context"

	"go-vents/internal/domain/services"

	status "google.golang.org/grpc/status"
)

func (s *Server) AllByNameAndSource(ctx context.Context, in *AllEventsByNameAndSourceRequest) (*AllEventsByNameAndSourceResponse, error) {
	var request = services.EventAllRequest{
		Name:   in.Name,
		Source: in.Source,
	}
	var service = services.NewEventAllService(s.repository)
	response, err := service.Do(ctx, request)
	if err != nil {
		return nil, status.Error(s.GetGRPCCode(err), err.Error())
	}

	return &AllEventsByNameAndSourceResponse{
		Events: s.ToProtos(response.Data),
	}, nil
}
