package userTicket

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r userTicket) GetOneByCode(ctx context.Context, orm *gorm.DB, req entity.UserTicket) (entity.UserTicket, error) {
	var resp entity.UserTicket
	if err := orm.WithContext(ctx).
		Model(&entity.UserTicket{}).
		Where("code = ?", req.Code).First(&resp).Error; err != nil {
		return resp, err
	}

	return resp, nil
}

type IGetOneByCode interface {
	GetOneByCode(ctx context.Context, orm *gorm.DB, req entity.UserTicket) (entity.UserTicket, error)
}
