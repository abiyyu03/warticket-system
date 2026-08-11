package ticket

import (
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"time"
)

type (
	InitOrderRequest struct {
		Date     string `json:"date"`
		EventID  int64  `json:"event_id"`
		Quantity int64  `json:"quantity"`
	}

	InitOrderResponse struct {
		Date     string `json:"date"`
		EventID  int64  `json:"event_id"`
		Quantity int64  `json:"quantity"`
	}
)

func (r InitOrderRequest) ToObEvent(parsedTime time.Time) entity.Event {
	return entity.Event{
		ID:        r.EventID,
		StartDate: parsedTime,
	}
}

func (r InitOrderRequest) ToObGetCache(userId int64) entity.CacheInitOrderRequest {
	return entity.CacheInitOrderRequest{
		UserID: userId,
	}
}

func (r InitOrderRequest) ToObSetCache(userId int64) entity.CacheInitOrderRequest {
	return entity.CacheInitOrderRequest{
		EventID:  r.EventID,
		Date:     r.Date,
		Quantity: r.Quantity,
		UserID:   userId,
	}
}
