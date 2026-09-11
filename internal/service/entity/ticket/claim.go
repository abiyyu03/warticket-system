package ticket

import ucEvent "go-projects/hexagonal-example/internal/service/entity/event"

type (
	PurchaseRequest struct {
		TxID string `json:"tx_id"` 
		Email   string               `json:"email"`
		Answers []ucEvent.AnswerInput `json:"answers"`
	}

	PurchaseResponse struct {
		Status string `json:"status"`
	}
)