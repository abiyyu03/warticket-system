package event

import (
	"context"
	obEntity "go-projects/hexagonal-example/internal/adapter/outbound/entity"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/event"
	"time"
)

func (s *service) CreateEvent(ctx context.Context, request ucEntity.CreateEventRequest) error {
	var (
		orm = s.repository.DB
	)
	start, err := time.Parse(time.RFC3339, request.StartDate)
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, request.EndDate)
	if err != nil {
		return err
	}

	// create event to db. event dilewatkan sebagai pointer supaya ID hasil
	// auto-increment terisi balik dan bisa dipakai sebagai key counter Redis.
	event := request.ToObEntity(start, end)
	if err := s.Repository.Event.Create(ctx, orm, &event); err != nil {
		return err
	}

	// set ticket quota value in redis. counter inilah yang di-DECR setiap
	// pembelian, jadi diisi sekali di sini dengan kuota awal event.
	err = s.Cache.Ticket.SetTicketQuota(ctx, obEntity.SetTicketQuotaRequest{
		EventID: event.ID,
		Quota:   request.Quota,
	})
	if err != nil {
		return err
	}

	return nil
}

type ICreateEvent interface {
	CreateEvent(ctx context.Context, request ucEntity.CreateEventRequest) error
}
