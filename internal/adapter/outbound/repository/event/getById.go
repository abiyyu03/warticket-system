package event

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

// GetByID mengambil event berdasarkan id tanpa validasi tanggal. Dipakai alur
// yang tidak lagi meminta tanggal dari client (date diisi server-side).
func (r event) GetByID(ctx context.Context, orm *gorm.DB, id int64) (entity.Event, error) {
	var event entity.Event
	err := orm.WithContext(ctx).Where("id = ?", id).First(&event).Error
	if err != nil {
		return event, err
	}
	return event, nil
}

type IGetByID interface {
	GetByID(ctx context.Context, orm *gorm.DB, id int64) (entity.Event, error)
}
