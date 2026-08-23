package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/event"

type (
	GetOneEventRequest struct {
		ID int64 `param:"id"`
	}
)

func (r GetOneEventRequest) ToUcEntity() ucEntity.GetOneEventRequest {
	return ucEntity.GetOneEventRequest{
		ID: r.ID,
	}
}
