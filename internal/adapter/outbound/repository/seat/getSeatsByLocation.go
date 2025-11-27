package seat

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r seat) GetSeatsByLocation(ctx context.Context, orm *gorm.DB, seat entity.Seat) ([]entity.Seat, error) {
	var result []entity.Seat
	if err := orm.WithContext(ctx).
		Where("location_id = ?", seat.LocationID).
		Where("available = ?", true).
		Find(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

type IGetSeatsByLocation interface {
	GetSeatsByLocation(ctx context.Context, orm *gorm.DB, seat entity.Seat) ([]entity.Seat, error)
}
