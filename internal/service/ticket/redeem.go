package ticket

import (
	"context"
	"errors"
	obEntity "go-projects/hexagonal-example/internal/adapter/outbound/entity"
	entity "go-projects/hexagonal-example/internal/service/entity/ticket"
	"go-projects/hexagonal-example/pkg/constants"

	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s service) Redeem(ctx context.Context, request entity.RedeemRequest) error {
	var orm = s.repository.DB

	trx := orm.Begin()
	defer func() {
		if r := recover(); r != nil {
			trx.Rollback()
		}
	}()

	// check if ticket code is valid
	userTicket, err := s.Repository.UserTicket.GetOneByCode(ctx, trx, obEntity.UserTicket{Code: request.TicketCode})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		trx.Rollback()
		return err
	}

	// check if ticket is already redeemed
	if userTicket.Status == constants.TicketStatusRedeemed {
		trx.Rollback()
		return errors.New("ticket already redeemed")
	}

	// check if ticket are not expired
	_, err = s.Repository.Event.GetOneByCode(ctx, trx, userTicket.Code)
	if err != nil {
		trx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("ticket not found", zap.Error(err))
			return errors.New("ticket not found")
		}
		return err
	}

	// redeem the ticket by updating its status and redeemAt timestamp
	err = s.Repository.UserTicket.RedeemTicket(ctx, trx, userTicket.Code)
	if err != nil {
		trx.Rollback()
		log.Error("Status update for ticket failed", zap.Error(err))
		return err
	}

	trx.Commit()

	// publish email notification to the user

	return nil
}

type IRedeem interface {
	Redeem(ctx context.Context, request entity.RedeemRequest) error
}
