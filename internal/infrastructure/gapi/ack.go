package gapi

import (
	context "context"

	"go-vents/internal/domain/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AckMessage(ctx context.Context, req *AckMessageRequest) (*AckMessageResponse, error) {
	id, err := types.NewSharedId(req.QueueId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.repository.AckMessage(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error ack: %v", err)
	}

	return &AckMessageResponse{
		Success: true,
	}, nil
}
