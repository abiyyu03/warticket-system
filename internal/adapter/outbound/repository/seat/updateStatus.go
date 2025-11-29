package seat

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r seat) UpdateStatus(ctx context.Context, orm *gorm.DB, seat entity.Seat) error {
	if err := orm.WithContext(ctx).
		Model(&entity.Seat{}).
		Where("id = ?", seat.ID).
		Update("available", seat.Available).Error; err != nil {
		return err
	}

	return nil
}

type IUpdateStatus interface {
	UpdateStatus(ctx context.Context, orm *gorm.DB, seat entity.Seat) error
}
