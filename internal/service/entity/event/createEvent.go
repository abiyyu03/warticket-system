package event

import (
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"mime/multipart"
	"time"
)

type (
	CreateEventRequest struct {
		Name        string
		Description string
		ImageFile   *multipart.FileHeader
		Price       float64
		Quota       int64
		StartDate   string
		EndDate     string
	}
)

func (r CreateEventRequest) ToObEntity(parsedStart time.Time, parsedEnd time.Time) entity.Event {
	return entity.Event{
		Name:        r.Name,
		Description: r.Description,
		ImageFile:   r.ImageFile.Filename,
		Price:       r.Price,
		Quota:       r.Quota,
		StartDate:   parsedStart,
		EndDate:     parsedEnd,
	}
}
