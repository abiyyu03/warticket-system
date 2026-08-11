package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

type (
	PurchaseRequest struct {
		EventID int64 `json:"event_id"`
	}
)

func (r PurchaseRequest) ToUcEntity() ucEntity.PurchaseRequest {
	return ucEntity.PurchaseRequest{
		EventID: r.EventID,
	}
}
