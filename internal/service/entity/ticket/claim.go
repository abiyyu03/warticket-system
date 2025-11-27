package ticket

import (
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"time"
)

type (
	ClaimTicketRequest struct {
		ChairCode *string `json:"chair_code"` // book for vip only
		Date      string  `json:"date"`
		EventID   int64   `json:"event_id"`
	}

	ClaimTicketResponse struct {
		Status string `json:"status"` //
	}
)

func (r ClaimTicketRequest) ToObSeat() entity.Seat {
	return entity.Seat{
		Code: *r.ChairCode,
	}
}
func (r ClaimTicketRequest) ToObSeatLocation(locationId int64) entity.Seat {
	return entity.Seat{
		LocationID: locationId,
	}
}

func (r ClaimTicketRequest) ToObEvent() entity.Event {
	return entity.Event{
		ID: r.EventID,
	}
}

func (r ClaimTicketRequest) ToObEventDetail(ticketType string) entity.EventDetail {
	return entity.EventDetail{
		EventID:    r.EventID,
		TicketType: ticketType,
	}
}

func (r ClaimTicketRequest) ToObTicket(userId, seatId, eventDetailId int64, code string, parsedTime time.Time) entity.UserEventTicket {
	return entity.UserEventTicket{
		UserID:        userId,
		SeatID:        seatId,
		EventDetailID: eventDetailId,
		TicketCode:    code,
		Date:          parsedTime,
	}
}
