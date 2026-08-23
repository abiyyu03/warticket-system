package event

import (
	"context"
	"errors"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/event"
	"strconv"
)

// SubmitRegistration menyimpan jawaban formulir peserta untuk sebuah event.
// Jawaban divalidasi terhadap definisi field; satu peserta hanya boleh mendaftar
// sekali per event.
func (s *service) SubmitRegistration(ctx context.Context, request ucEntity.SubmitRegistrationRequest) error {
	orm := s.repository.DB

	userIdCtx, _ := ctx.Value("x-user-id").(string)
	userId, _ := strconv.ParseInt(userIdCtx, 10, 64)
	request.UserID = userId

	// event harus punya formulir; kalau tidak, tidak ada yang perlu didaftarkan.
	fields, err := s.Repository.Event.GetFormFieldsByEvent(ctx, orm, request.EventID)
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return errors.New("event ini tidak memiliki formulir pendaftaran")
	}

	if err := request.Validate(fields); err != nil {
		return err
	}

	// cegah pendaftaran ganda (backstop: unique constraint user+event).
	exists, err := s.Repository.UserRegistration.ExistsByUserEvent(ctx, orm, userId, request.EventID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("kamu sudah terdaftar pada event ini")
	}

	reg := request.ToObEntity()
	return s.Repository.UserRegistration.Create(ctx, orm, &reg)
}

type ISubmitRegistration interface {
	SubmitRegistration(ctx context.Context, request ucEntity.SubmitRegistrationRequest) error
}
