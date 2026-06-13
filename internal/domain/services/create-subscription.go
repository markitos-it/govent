package services

import (
	"context"

	"go-vents/internal/domain/shared"
	"go-vents/internal/domain/types"
)

type SubscriptionCreateRequest struct {
	SubscriberName string
	Source         string
	Event          string
}

type SubscriptionCreateResponse struct {
	Id             string `json:"id"`
	SubscriberName string `json:"subscriber_name"`
	Source         string `json:"source"`
	Event          string `json:"event"`
}

type SubscriptionCreateService struct {
	Repository types.EventRepository
}

func NewSubscriptionCreateService(repository types.EventRepository) SubscriptionCreateService {
	return SubscriptionCreateService{
		Repository: repository,
	}
}

func (s SubscriptionCreateService) Do(ctx context.Context, request SubscriptionCreateRequest) (*SubscriptionCreateResponse, error) {

	subscription, err := types.NewSubscription(shared.UUIDv4(), request.SubscriberName, request.Event, request.Source)
	if err != nil {
		return nil, err
	}

	if err := s.Repository.CreateSubscription(ctx, subscription); err != nil {
		return nil, err
	}

	return &SubscriptionCreateResponse{
		Id:             subscription.Id,
		SubscriberName: subscription.SubscriberName,
		Source:         subscription.Source,
		Event:          subscription.EventName,
	}, nil
}
