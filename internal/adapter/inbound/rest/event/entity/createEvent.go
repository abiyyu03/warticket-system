package entity

import (
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
