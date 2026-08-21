package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/event"

type (
	SubmitRegistrationRequest struct {
		Answers []AnswerPayload `json:"answers"`
	}

	AnswerPayload struct {
		FieldID int64    `json:"field_id"`
		// value selalu array supaya seragam: text/select berisi satu, checkbox banyak.
		Value []string `json:"value"`
	}
)

func (r SubmitRegistrationRequest) ToUcEntity(eventID int64) ucEntity.SubmitRegistrationRequest {
	answers := make([]ucEntity.AnswerInput, 0, len(r.Answers))
	for _, a := range r.Answers {
		answers = append(answers, ucEntity.AnswerInput{
			FieldID: a.FieldID,
			Value:   a.Value,
		})
	}
	return ucEntity.SubmitRegistrationRequest{
		EventID: eventID,
		Answers: answers,
	}
}
