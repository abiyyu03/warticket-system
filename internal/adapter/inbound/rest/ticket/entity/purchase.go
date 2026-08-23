package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"

type (
	PurchaseRequest struct {
		TxID string `json:"tx_id"`
	}
)

func (r PurchaseRequest) ToUcEntity() ucEntity.PurchaseRequest {
	return ucEntity.PurchaseRequest{
		TxID: r.TxID,
	}
}
