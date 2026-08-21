package event

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

// GetFormFieldsByEvent mengambil definisi field formulir sebuah event, urut
// berdasarkan position lalu id.
func (r event) GetFormFieldsByEvent(ctx context.Context, orm *gorm.DB, eventID int64) ([]entity.EventFormField, error) {
	var fields []entity.EventFormField
	err := orm.WithContext(ctx).
		Where("event_id = ?", eventID).
		Order("position ASC, id ASC").
		Find(&fields).Error
	if err != nil {
		return nil, err
	}
	return fields, nil
}

type IGetFormFields interface {
	GetFormFieldsByEvent(ctx context.Context, orm *gorm.DB, eventID int64) ([]entity.EventFormField, error)
}
