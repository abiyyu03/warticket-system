package userTicket

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r userTicket) RedeemTicket(ctx context.Context, orm *gorm.DB, code string) error {
	var ticket entity.UserTicket

	// lock
	if err := orm.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("code = ?", code).
		First(&ticket).Error; err != nil {
		return err
	}

	// update
	return orm.WithContext(ctx).
		Model(&entity.UserTicket{}).
		Where("id = ?", ticket.ID).
		Updates(map[string]interface{}{
			"status":      entity.UserTicketStatusRedeemed,
			"redeemed_at": time.Now(),
		}).Error
}

type IRedeemTicket interface {
	RedeemTicket(ctx context.Context, orm *gorm.DB, code string) error
}
