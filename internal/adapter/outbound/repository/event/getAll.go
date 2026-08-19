package event

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r event) GetAll(ctx context.Context, orm *gorm.DB, filter entity.Event) ([]entity.Event, error) {
	var events []entity.Event

	query := orm.WithContext(ctx).Model(&entity.Event{})
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Order("start_date DESC").Find(&events).Error; err != nil {
		return events, err
	}

	return events, nil
}

type IGetAll interface {
	GetAll(ctx context.Context, orm *gorm.DB, filter entity.Event) ([]entity.Event, error)
}
