package gapi

import (
	context "context"

	"go-vents/internal/domain/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DEFAULT_LIMIT_PULL = 10

func (s *Server) PullMessages(ctx context.Context, req *PullMessagesRequest) (*PullMessagesResponse, error) {
	event, err := types.NewEventName(req.GetEventName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	source, err := types.NewEventSource(req.GetSource())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dbMessages, err := s.repository.PullMessages(ctx, event, source)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error en la base de datos al hacer pull: %v", err)
	}

	var protoMessages []*QueueMessage
	for _, msg := range dbMessages {
		protoMessages = append(protoMessages, &QueueMessage{
			QueueId: msg.Id,
			EventId: msg.EventId,
			Name:    msg.SubscriberName,
			Source:  "fakesource",
			Payload: "fakepayload",
		})
	}

	return &PullMessagesResponse{
		Messages: protoMessages,
	}, nil
}
