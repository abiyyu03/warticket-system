package seat

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r seat) GetOne(ctx context.Context, orm *gorm.DB, seat entity.Seat) (entity.Seat, error) {
	var result entity.Seat
	if err := orm.WithContext(ctx).
		Where("code = ?", seat.Code).
		Where("available = ?", seat.Available).
		First(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

type IGetOne interface {
	GetOne(ctx context.Context, orm *gorm.DB, seat entity.Seat) (entity.Seat, error)
}
