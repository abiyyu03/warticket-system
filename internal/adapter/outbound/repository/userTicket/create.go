package userTicket

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r userTicket) Create(ctx context.Context, orm *gorm.DB, req entity.UserTicket) error {
	if err := orm.WithContext(ctx).Model(&entity.UserTicket{}).Create(&req).Error; err != nil {
		return err
	}

	return nil
}

type ICreate interface {
	Create(ctx context.Context, orm *gorm.DB, req entity.UserTicket) error
}
