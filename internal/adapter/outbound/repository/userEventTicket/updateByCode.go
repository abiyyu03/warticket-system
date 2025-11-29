package userEventTicket

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r userEventTicket) UpdateByCode(ctx context.Context, orm *gorm.DB, ticket entity.UserEventTicket) error {
	if err := orm.WithContext(ctx).
		Model(&entity.UserEventTicket{}).
		Where("ticket_code = ?", ticket.TicketCode).
		Updates(&ticket).Error; err != nil {
		return err
	}

	return nil
}

type IUpdateByCode interface {
	UpdateByCode(ctx context.Context, orm *gorm.DB, ticket entity.UserEventTicket) error
}
