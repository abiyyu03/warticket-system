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

	// validasi definisi formulir sebelum menyentuh DB.
	if err := request.ValidateFormFields(); err != nil {
		return err
	}

	// event + form field-nya dibuat atomik dalam satu transaksi.
	trx := orm.Begin()
	defer func() {
		if r := recover(); r != nil {
			trx.Rollback()
		}
	}()

	event := request.ToObEntity(start, end)
	if err := s.Repository.Event.Create(ctx, trx, &event); err != nil {
		trx.Rollback()
		return err
	}

	// form field opsional; id event baru dipakai sebagai foreign key.
	if fields := request.ToObFormFields(event.ID); len(fields) > 0 {
		if err := s.Repository.Event.CreateFormFields(ctx, trx, fields); err != nil {
			trx.Rollback()
			return err
		}
	}

	if err := trx.Commit().Error; err != nil {
		return err
	}

	// counter kuota di redis di-set setelah event tersimpan permanen.
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
