package event

import (
	"time"

	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

type (
	GetListEventRequest struct {
		Page   int    `query:"page"`
		Limit  int    `query:"limit"`
		Search string `query:"search"`
	}

	Event struct {
		ID          int64     `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description,omitempty"`
		StartDate   time.Time `json:"start_date"`
		EndDate     time.Time `json:"end_date"`
		Price       float64   `json:"price"`
		Quota       int64     `json:"quota"`
		ImageFile   string    `json:"image_file,omitempty"`
	}

	GetListEventResponse struct {
		Events []Event `json:"events"`
		Page   int     `json:"page"`
		Limit  int     `json:"limit"`
		Total  int64   `json:"total"`
	}
)

// ToObEntity memetakan filter list ke entity. Search dipakai untuk mencari nama event.
func (r GetListEventRequest) ToObEntity() entity.Event {
	return entity.Event{
		Name: r.Search,
	}
}
