package entity

import (
	ucEvent "go-projects/hexagonal-example/internal/service/entity/event"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/ticket"
)

type (
	PurchaseRequest struct {
		TxID    string          `json:"tx_id"`
		Email   string          `json:"email"`
		Answers []AnswerPayload `json:"answers"`
	}

	// AnswerPayload: satu jawaban formulir; value selalu array (text/select isi
	// satu, checkbox banyak).
	AnswerPayload struct {
		FieldID int64    `json:"field_id"`
		Value   []string `json:"value"`
	}
)

func (r PurchaseRequest) ToUcEntity() ucEntity.PurchaseRequest {
	answers := make([]ucEvent.AnswerInput, 0, len(r.Answers))
	for _, a := range r.Answers {
		answers = append(answers, ucEvent.AnswerInput{FieldID: a.FieldID, Value: a.Value})
	}
	return ucEntity.PurchaseRequest{
		TxID:    r.TxID,
		Email:   r.Email,
		Answers: answers,
	}
}
