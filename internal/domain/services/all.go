package services

import (
	"context"
	"go-vents/internal/domain/types"
)

type EventAllResponse struct {
	Data []*types.Event `json:"data"`
}

type EventAllRequest struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type EventAllService struct {
	Repository types.EventRepository
}

func NewEventAllService(repository types.EventRepository) EventAllService {
	return EventAllService{
		Repository: repository,
	}
}

func (s EventAllService) Do(ctx context.Context, request EventAllRequest) (*EventAllResponse, error) {
	eventName, err := types.NewName(request.Name)
	if err != nil {
		return nil, err
	}
	eventSource, err := types.NewSource(request.Source)
	if err != nil {
		return nil, err
	}

	events, err := s.Repository.AllByNameAndSource(ctx, eventName, eventSource)
	if err != nil {
		return nil, err
	}

	return &EventAllResponse{
		Data: events,
	}, nil
}
