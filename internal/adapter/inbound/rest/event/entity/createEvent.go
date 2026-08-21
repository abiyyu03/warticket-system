package entity

import (
	"encoding/json"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/event"
	"mime/multipart"
)

type (
	CreateEventRequest struct {
		Name        string                `form:"name"`
		Description string                `form:"description"`
		ImageFile   *multipart.FileHeader `form:"image_file"`
		Price       float64               `form:"price"`
		Quota       int64                 `form:"quota"`
		StartDate   string                `form:"start_date"`
		EndDate     string                `form:"end_date"`
		FormFields  string                `form:"form_fields"`
	}

	// FormFieldPayload adalah bentuk JSON satu field yang dikirim client.
	FormFieldPayload struct {
		Label     string   `json:"label"`
		FieldType string   `json:"field_type"`
		Required  bool     `json:"required"`
		Options   []string `json:"options"`
		Position  int      `json:"position"`
	}
)

func (r CreateEventRequest) ToUcEntity() ucEntity.CreateEventRequest {
	return ucEntity.CreateEventRequest{
		Name:        r.Name,
		Description: r.Description,
		ImageFile:   r.ImageFile,
		Price:       r.Price,
		Quota:       r.Quota,
		StartDate:   r.StartDate,
		EndDate:     r.EndDate,
	}
}

// ParseFormFields mem-parse string JSON form_fields menjadi input usecase.
// String kosong berarti event tanpa formulir kustom.
func (r CreateEventRequest) ParseFormFields() ([]ucEntity.FormFieldInput, error) {
	if r.FormFields == "" {
		return nil, nil
	}

	var payload []FormFieldPayload
	if err := json.Unmarshal([]byte(r.FormFields), &payload); err != nil {
		return nil, err
	}

	fields := make([]ucEntity.FormFieldInput, 0, len(payload))
	for _, p := range payload {
		fields = append(fields, ucEntity.FormFieldInput{
			Label:     p.Label,
			FieldType: p.FieldType,
			Required:  p.Required,
			Options:   p.Options,
			Position:  p.Position,
		})
	}
	return fields, nil
}
