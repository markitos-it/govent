package gapi

import (
	context "context"

	"go-vents/internal/domain/services"
	"go-vents/internal/domain/types"

	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *Server) GetEvent(ctx context.Context, in *GetEventRequest) (*GetEventResponse, error) {
	if _, err := types.NewSharedId(in.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	request := services.EventOneRequest{Id: in.Id}

	var service = services.NewEventOneService(s.repository)
	response, err := service.Do(ctx, request)
	if err != nil {
		return nil, status.Error(s.GetGRPCCode(err), err.Error())

	}

	return &GetEventResponse{
		Id:      response.Data.Id,
		Name:    response.Data.Name,
		Source:  response.Data.Source,
		Payload: response.Data.Payload,
	}, nil
}
