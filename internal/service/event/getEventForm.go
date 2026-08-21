package event

import (
	"context"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/event"
)

// GetEventForm mengembalikan definisi formulir pendaftaran sebuah event.
// Fields kosong berarti event tidak memakai formulir kustom.
func (s *service) GetEventForm(ctx context.Context, eventID int64) (ucEntity.GetEventFormResponse, error) {
	fields, err := s.Repository.Event.GetFormFieldsByEvent(ctx, s.repository.DB, eventID)
	if err != nil {
		return ucEntity.GetEventFormResponse{}, err
	}
	return ucEntity.ToGetEventFormResponse(eventID, fields), nil
}

type IGetEventForm interface {
	GetEventForm(ctx context.Context, eventID int64) (ucEntity.GetEventFormResponse, error)
}
