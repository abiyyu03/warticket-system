package userTicket

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r userTicket) ListByEventWithEmail(ctx context.Context, orm *gorm.DB, eventID int64) ([]entity.EventTicketExport, error) {
	var rows []entity.EventTicketExport
	err := orm.WithContext(ctx).
		Table("user_tickets AS t").
		Select("t.id, t.code, t.user_id, COALESCE(r.email, '') AS email, t.status, t.valid_until, t.created_at").
		Joins("LEFT JOIN user_registrations r ON r.user_id = t.user_id AND r.event_id = t.event_id").
		Where("t.event_id = ?", eventID).
		Order("t.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

type IListByEventWithEmail interface {
	ListByEventWithEmail(ctx context.Context, orm *gorm.DB, eventID int64) ([]entity.EventTicketExport, error)
}
