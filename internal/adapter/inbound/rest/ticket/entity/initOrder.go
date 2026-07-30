package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

type (
	InitOrderRequest struct {
		ChairCode []string `json:"chair_code"`
		Date      string   `json:"date"`
		EventID   int64    `json:"event_id"`
	}
)

func (r InitOrderRequest) ToUcEntity() ucEntity.InitOrderRequest {
	return ucEntity.InitOrderRequest{
		Date:      r.Date,
		EventID:   r.EventID,
		ChairCode: r.ChairCode,
	}
}
