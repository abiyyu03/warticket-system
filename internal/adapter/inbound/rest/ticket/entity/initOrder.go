package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

type (
	InitOrderRequest struct {
		Date     string `json:"date" validate:"required,datetime=2006-01-02"`
		EventID  int64  `json:"event_id" validate:"required"`
		Quantity int64  `json:"quantity" validate:"required,gte=1"`
	}
)

func (r InitOrderRequest) ToUcEntity() ucEntity.InitOrderRequest {
	return ucEntity.InitOrderRequest{
		Date:     r.Date,
		EventID:  r.EventID,
		Quantity: r.Quantity,
	}
}
