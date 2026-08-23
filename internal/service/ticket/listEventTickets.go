package ticket

import (
	"context"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
)

// ListEventTickets mengembalikan seluruh tiket sebuah event beserta email
// peserta (pengelolaan sisi author).
func (s service) ListEventTickets(ctx context.Context, eventID int64) (ucEntity.ListEventTicketsResponse, error) {
	rows, err := s.Repository.UserTicket.ListByEventWithEmail(ctx, s.repository.DB, eventID)
	if err != nil {
		return ucEntity.ListEventTicketsResponse{}, err
	}
	return ucEntity.ToListEventTicketsResponse(eventID, rows), nil
}

type IListEventTickets interface {
	ListEventTickets(ctx context.Context, eventID int64) (ucEntity.ListEventTicketsResponse, error)
}
