package services

import (
	"context"

	"govent/internal/domain/shared"
	"govent/internal/domain/types"
)

type EventCreateRequest struct {
	Name    string
	Source  string
	Payload string
}

type EventCreateResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Source  string `json:"source"`
	Payload string `json:"payload"`
}

type EventCreateService struct {
	Repository types.EventRepository
}

func NewEventCreateService(repository types.EventRepository) EventCreateService {
	return EventCreateService{
		Repository: repository,
	}
}

func (s EventCreateService) Do(ctx context.Context, request EventCreateRequest) (*EventCreateResponse, error) {

	event, err := types.NewEvent(shared.UUIDv4(), request.Name, request.Source, request.Payload)
	if err != nil {
		return nil, err
	}

	if err := s.Repository.Create(ctx, event); err != nil {
		return nil, err
	}

	return &EventCreateResponse{
		Id:      event.Id,
		Name:    event.Name,
		Source:  event.Source,
		Payload: event.Payload,
	}, nil
}
