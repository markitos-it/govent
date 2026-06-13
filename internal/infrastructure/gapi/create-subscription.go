package gapi

import (
	context "context"

	"go-vents/internal/domain/services"

	"google.golang.org/grpc/status"
)

func (s *Server) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*CreateSubscriptionResponse, error) {
	var request = services.SubscriptionCreateRequest{
		SubscriberName: req.SubscriberName,
		Source:         req.Source,
		Event:          req.EventName,
	}

	var service = services.NewSubscriptionCreateService(s.repository)
	entity, err := service.Do(ctx, request)
	if err != nil {
		return nil, status.Error(s.GetGRPCCode(err), err.Error())
	}

	return &CreateSubscriptionResponse{
		Success: true,
		Message: "Subscription created successfully id: " + entity.Id,
	}, nil
}
