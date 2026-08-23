package event

import (
	"context"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/event"
)

func (s *service) GetOneEvent(ctx context.Context, request ucEntity.GetOneEventRequest) (ucEntity.GetOneEventResponse, error) {
	var (
		orm      = s.repository.DB
		response ucEntity.GetOneEventResponse
	)

	event, err := s.Repository.Event.GetOneById(ctx, orm, request.ToObEntity())
	if err != nil {
		return response, err
	}

	response.Event = ucEntity.Event{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
	}

	return response, nil
}

type IGetOneEvent interface {
	GetOneEvent(ctx context.Context, request ucEntity.GetOneEventRequest) (ucEntity.GetOneEventResponse, error)
}
