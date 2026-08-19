package event

import (
	"context"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/event"
)

func (s *service) GetListEvent(ctx context.Context, request ucEntity.GetListEventRequest) (ucEntity.GetListEventResponse, error) {
	var (
		orm      = s.repository.DB
		response ucEntity.GetListEventResponse
	)

	list, err := s.Repository.Event.GetAll(ctx, orm, request.ToObEntity())
	if err != nil {
		return response, err
	}

	returnedList := make([]ucEntity.Event, len(list))
	for i, event := range list {
		returnedList[i] = ucEntity.Event{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
		}
	}
	response.Events = returnedList

	return response, nil
}

type IGetListEvent interface {
	GetListEvent(ctx context.Context, request ucEntity.GetListEventRequest) (ucEntity.GetListEventResponse, error)
}
