package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

type (
	InitOrderRequest struct {
		Date     string `json:"date"`
		EventID  int64  `json:"event_id"`
		Quantity int64  `json:"quantity"`
	}
)

func (r InitOrderRequest) ToUcEntity() ucEntity.InitOrderRequest {
	return ucEntity.InitOrderRequest{
		Date:     r.Date,
		EventID:  r.EventID,
		Quantity: r.Quantity,
	}
}
