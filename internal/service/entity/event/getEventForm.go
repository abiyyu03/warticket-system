package event

import "go-projects/hexagonal-example/internal/adapter/outbound/entity"

type (
	EventFormFieldResponse struct {
		ID        int64    `json:"id"`
		Label     string   `json:"label"`
		FieldType string   `json:"field_type"`
		Required  bool     `json:"required"`
		Options   []string `json:"options,omitempty"`
		Position  int      `json:"position"`
	}

	GetEventFormResponse struct {
		EventID int64                    `json:"event_id"`
		Fields  []EventFormFieldResponse `json:"fields"`
	}
)

// ToGetEventFormResponse memetakan definisi field outbound ke response.
func ToGetEventFormResponse(eventID int64, fields []entity.EventFormField) GetEventFormResponse {
	out := make([]EventFormFieldResponse, 0, len(fields))
	for _, f := range fields {
		out = append(out, EventFormFieldResponse{
			ID:        f.ID,
			Label:     f.Label,
			FieldType: f.FieldType,
			Required:  f.Required,
			Options:   []string(f.Options),
			Position:  f.Position,
		})
	}
	return GetEventFormResponse{EventID: eventID, Fields: out}
}
