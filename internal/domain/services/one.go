package services

import (
	"context"

	"govent/internal/domain/types"
)

type EventOneRequest struct {
	Id string `json:"id"`
}

type EventOneResponse struct {
	Data *types.Event `json:"data"`
}

type EventOneService struct {
	Repository types.EventRepository
}

func NewEventOneService(repository types.EventRepository) EventOneService {
	return EventOneService{
		Repository: repository,
	}
}

func (s EventOneService) Do(ctx context.Context, request EventOneRequest) (*EventOneResponse, error) {
	securedId, err := types.NewEventId(request.Id)
	if err != nil {
		return nil, err
	}

	event, err := s.Repository.One(ctx, securedId)
	if err != nil {
		return nil, err
	}

	return &EventOneResponse{event}, nil
}
