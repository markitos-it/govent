package gapi

import (
	context "context"

	"go-vents/internal/domain/services"
	"go-vents/internal/domain/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteEvent(ctx context.Context, in *DeleteEventRequest) (*DeleteEventResponse, error) {
	if _, err := types.NewSharedId(in.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := services.EventDeleteRequest{Id: in.Id}

	var service = services.NewEventDeleteService(s.repository)
	if err := service.Do(ctx, request); err != nil {
		return nil, status.Error(s.GetGRPCCode(err), err.Error())
	}

	return &DeleteEventResponse{
		Deleted: request.Id,
	}, nil
}
