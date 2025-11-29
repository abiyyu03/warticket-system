package seat

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r seat) UpdateByCode(ctx context.Context, orm *gorm.DB, seat entity.Seat) error {
	if err := orm.WithContext(ctx).
		Where("code = ?", seat.Code).
		Updates(&seat).Error; err != nil {
		return err
	}

	return nil
}

type IUpdateByCode interface {
	UpdateByCode(ctx context.Context, orm *gorm.DB, seat entity.Seat) error
}
