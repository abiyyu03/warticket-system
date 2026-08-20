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

func (r CreateEventRequest) ToObEntity(parsedStart, parsedEnd time.Time) entity.Event {
	var imageFile string
	if r.ImageFile != nil {
		imageFile = r.ImageFile.Filename
	}

	return entity.Event{
		Name:        r.Name,
		Description: r.Description,
		ImageFile:   imageFile,
		Price:       r.Price,
		Quota:       r.Quota,
		QuotaRemaining: r.Quota,
		StartDate:      parsedStart,
		EndDate:        parsedEnd,
	}
}
