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
		Price    int64  `json:"price"`
	}

	InitOrderResponse struct {
		Date     string `json:"date"`
		EventID  int64  `json:"event_id"`
		Quantity int64  `json:"quantity"`
		Price    int64  `json:"price"`
	}
)

func (r InitOrderRequest) ToObEvent(parsedDate time.Time) entity.Event {
	return entity.Event{
		ID:        r.EventID,
		StartDate: parsedDate,
	}
}

func (r InitOrderRequest) ToObGetCache(userId int64) entity.CacheInitOrderRequest {
	return entity.CacheInitOrderRequest{
		UserID: userId,
	}
}

// ToObSetCache menyusun payload cache. price berasal dari event (server-side),
// bukan dari request.
func (r InitOrderRequest) ToObSetCache(userId, price int64) entity.CacheInitOrderRequest {
	return entity.CacheInitOrderRequest{
		EventID:  r.EventID,
		Date:     r.Date,
		Quantity: r.Quantity,
		UserID:   userId,
		Price:    price,
	}
}
