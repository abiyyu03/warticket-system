package event

import (
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

type (
	GetOneEventRequest struct {
		ID int64 `json:"id"`
	}

	GetOneEventResponse struct {
		Event Event `json:"event"`
	}
)

// ToObEntity memetakan filter list ke entity. Search dipakai untuk mencari nama event.
func (r GetOneEventRequest) ToObEntity() entity.Event {
	return entity.Event{
		ID: r.ID,
	}
}
