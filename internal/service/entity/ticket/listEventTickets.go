package ticket

import (
	"time"

	obEntity "go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

type (
	EventTicket struct {
		ID         int64     `json:"id"`
		Code       string    `json:"code"`
		UserID     int64     `json:"user_id"`
		Email      string    `json:"email,omitempty"`
		Status     string    `json:"status"`
		ValidUntil time.Time `json:"valid_until"`
		CreatedAt  time.Time `json:"created_at"`
	}

	ListEventTicketsResponse struct {
		EventID int64         `json:"event_id"`
		Total   int           `json:"total"`
		Tickets []EventTicket `json:"tickets"`
	}
)

// ToListEventTicketsResponse memetakan proyeksi tiket+email ke response ringkas
// untuk pengelolaan sisi author.
func ToListEventTicketsResponse(eventID int64, rows []obEntity.EventTicketExport) ListEventTicketsResponse {
	out := make([]EventTicket, 0, len(rows))
	for _, r := range rows {
		out = append(out, EventTicket{
			ID:         r.ID,
			Code:       r.Code,
			UserID:     r.UserID,
			Email:      r.Email,
			Status:     r.Status,
			ValidUntil: r.ValidUntil,
			CreatedAt:  r.CreatedAt,
		})
	}
	return ListEventTicketsResponse{
		EventID: eventID,
		Total:   len(out),
		Tickets: out,
	}
}
