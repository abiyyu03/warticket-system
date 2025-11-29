package ticket

import (
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"time"
)

type (
	ClaimTicketRequest struct{}

	ClaimAndPurchaseTicketResponse struct {
		Status string `json:"status"` //
	}
)

func (r ClaimTicketRequest) ToObSeat(chairCode string) entity.Seat {
	return entity.Seat{
		Code:      chairCode,
		Available: true,
	}
}

func (r ClaimTicketRequest) ToObSeatUnavailable(seatId int64) entity.Seat {
	return entity.Seat{
		ID:        seatId,
		Available: false,
	}
}

func (r ClaimTicketRequest) ToObSeatLocation(locationId int64) entity.Seat {
	return entity.Seat{
		LocationID: locationId,
	}
}

func (r ClaimTicketRequest) ToObEvent(eventId int64) entity.Event {
	return entity.Event{
		ID: eventId,
	}
}

func (r ClaimTicketRequest) ToObGetCache(userId int64) entity.CacheInitOrderRequest {
	return entity.CacheInitOrderRequest{
		UserID: userId,
	}
}

func (r ClaimTicketRequest) ToObEventDetail(eventId int64, ticketType string) entity.EventDetail {
	return entity.EventDetail{
		EventID:    eventId,
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
