package entity

import ucEntity "go-projects/hexagonal-example/internal/service/entity/event"

type (
	GetListEventRequest struct {
		Page   int    `query:"page"`
		Limit  int    `query:"limit"`
		Search string `query:"search"`
	}
)

func (r GetListEventRequest) ToUcEntity() ucEntity.GetListEventRequest {
	return ucEntity.GetListEventRequest{
		Page:   r.Page,
		Limit:  r.Limit,
		Search: r.Search,
	}
}
