package event

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r event) GetOneByCode(ctx context.Context, orm *gorm.DB, code string) (entity.Event, error) {
	var eventResult entity.Event
	err := orm.WithContext(ctx).
		Where("code = ?", code).
		First(&eventResult).Error
	if err != nil {
		return eventResult, err
	}

	return eventResult, nil
}

type IGetOneByCode interface {
	GetOneByCode(ctx context.Context, orm *gorm.DB, code string) (entity.Event, error)
}
