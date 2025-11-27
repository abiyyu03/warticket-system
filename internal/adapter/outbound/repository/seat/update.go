package seat

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r seat) Update(ctx context.Context, orm *gorm.DB, seat entity.Seat) error {
	if err := orm.WithContext(ctx).
		Where("code = ?", seat.Code).
		Updates(&seat).Error; err != nil {
		return err
	}

	return nil
}

type IUpdate interface {
	Update(ctx context.Context, orm *gorm.DB, seat entity.Seat) error
}
