package userEvent

import (
	"context"
	"go-projects/hexagonal-example/internal/adapter/outbound/entity"
)

func (r user) Create(ctx context.Context, req entity.UserEventTicket) error {
	if err := r.Package.DB.WithContext(ctx).Create(&req).Error; err != nil {
		return err
	}

	return nil
}

type ICreate interface {
	Create(ctx context.Context, req entity.UserEventTicket) error
}
