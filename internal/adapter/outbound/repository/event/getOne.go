package event

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r event) GetOneById(ctx context.Context, orm *gorm.DB, event entity.Event) (entity.Event, error) {
	var result entity.Event
	query := orm.WithContext(ctx).
		Where("id = ?", event.ID).
		First(&result)

	if event.EndDate != nil {
		query = query.Where("start_date <=", event.StartDate).Where("end_date >= ", event.StartDate)
	} else {
		query = query.Where("start_date <=", event.StartDate) // the event only one day
	}

	if err := query.Error; err != nil {
		return result, err

	}
	return result, nil
}

type IGetOneById interface {
	GetOneById(ctx context.Context, orm *gorm.DB, event entity.Event) (entity.Event, error)
}
