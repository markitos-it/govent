package services

import (
	"context"

	"go-vents/internal/domain/types"
)

type EventDeleteRequest struct {
	Id string `json:"id"`
}

type EventDeleteService struct {
	Repository types.EventRepository
}

func NewEventDeleteService(repository types.EventRepository) EventDeleteService {
	return EventDeleteService{
		Repository: repository,
	}
}

func (s EventDeleteService) Do(ctx context.Context, request EventDeleteRequest) error {
	securedId, err := types.NewSharedId(request.Id)
	if err != nil {
		return err
	}

	if err := s.Repository.Delete(ctx, securedId); err != nil {
		return err
	}

	return nil
}
