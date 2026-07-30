package event

import (
	"context"
	"errors"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
	"time"

	"gorm.io/gorm"
)

func (r event) GetOneById(ctx context.Context, orm *gorm.DB, event entity.Event) (entity.Event, error) {
	var eventResult entity.Event
	err := orm.WithContext(ctx).
		Where("id = ?", event.ID).
		First(&eventResult).Error
	if err != nil {
		return eventResult, err
	}

	userDate := event.StartDate.Truncate(24 * time.Hour)
	startDate := eventResult.StartDate.Truncate(24 * time.Hour)

	// EndDate value nol (IsZero) berarti event satu hari.
	// Single-day event
	if eventResult.EndDate.IsZero() {
		if userDate.Equal(startDate) {
			return eventResult, nil
		}
		return eventResult, errors.New("input date out of range for single-day event")
	}

	// Multi-day event
	endDate := eventResult.EndDate.Truncate(24 * time.Hour)
	if userDate.Before(startDate) || userDate.After(endDate) {
		return eventResult, errors.New("input date outside event range")
	}

	return eventResult, nil
}

type IGetOneById interface {
	GetOneById(ctx context.Context, orm *gorm.DB, event entity.Event) (entity.Event, error)
}
