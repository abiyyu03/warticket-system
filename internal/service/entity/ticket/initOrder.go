package ticket

import (
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"time"
)

type (
	InitOrderRequest struct {
		ChairCode []string `json:"chair_code"`
		Date      string   `json:"date"`
		EventID   int64    `json:"event_id"`
	}

	InitOrderResponse struct {
		Status string `json:"status"` //
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

func (r InitOrderRequest) ToObSetCache(userId int64) entity.CacheInitOrderResponse {
	return entity.CacheInitOrderResponse{
		ChairCode: r.ChairCode,
		EventID:   r.EventID,
		Date:      r.Date,
		UserID:    userId,
	}
}
