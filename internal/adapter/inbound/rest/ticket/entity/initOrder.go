package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

type (
	InitOrderRequest struct {
		EventID  int64 `json:"event_id" validate:"required"`
		Quantity int64 `json:"quantity" validate:"required,gte=1"`
	}
)

func (r InitOrderRequest) ToUcEntity() ucEntity.InitOrderRequest {
	return ucEntity.InitOrderRequest{
		EventID:  r.EventID,
		Quantity: r.Quantity,
	}
}
