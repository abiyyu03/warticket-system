package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/event"

type (
	CreateEventRequest struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		ImageFile   string  `json:"image_file"`
		Price       float64 `json:"price"`
		Quota       int64   `json:"quota"`
		StartDate   string  `json:"start_date"`
		EndDate     string  `json:"end_date"`
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
