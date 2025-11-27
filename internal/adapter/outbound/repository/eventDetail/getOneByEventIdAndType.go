package eventDetail

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"

	"gorm.io/gorm"
)

func (r eventDetail) GetOneByEventIdAndType(ctx context.Context, orm *gorm.DB, eventDetail entity.EventDetail) (entity.EventDetail, error) {
	var result entity.EventDetail
	if err := orm.WithContext(ctx).
		Where("event_id = ?", eventDetail.EventID).
		Where("ticket_type = ?", eventDetail.TicketType).
		First(&result).Error; err != nil {
		return result, err

	}
	return result, nil
}

type IGetOneByEventIdAndType interface {
	GetOneByEventIdAndType(ctx context.Context, orm *gorm.DB, eventDetail entity.EventDetail) (entity.EventDetail, error)
}
