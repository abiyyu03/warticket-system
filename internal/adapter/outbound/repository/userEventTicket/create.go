package userEventTicket

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r userEventTicket) Create(ctx context.Context, orm *gorm.DB, req entity.UserEventTicket) error {
	if err := orm.WithContext(ctx).Model(&entity.UserEventTicket{}).Create(&req).Error; err != nil {
		return err
	}

	return nil
}

type ICreate interface {
	Create(ctx context.Context, orm *gorm.DB, req entity.UserEventTicket) error
}
