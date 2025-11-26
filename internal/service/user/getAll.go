package user

import (
	"context"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/user"
)

func (s service) GetAll(ctx context.Context) ([]ucEntity.User, error) {
	var response []ucEntity.User

	cachedUser, err := s.Cache.User.GetAllUser(ctx)
	if err == nil {
		for _, value := range cachedUser {
			response = append(response, ucEntity.User{
				ID:    value.ID,
				Email: value.Email,
				Name:  value.Name,
			})
		}
		return response, nil
	}

	users, err := s.Repository.User.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	err = s.Cache.User.SetAllUser(ctx, users)
	if err != nil {
		return nil, err
	}

	for _, value := range users {
		response = append(response, ucEntity.User{
			ID:    value.ID,
			Email: value.Email,
			Name:  value.Name,
		})
	}

	return response, nil
}

type IGetAll interface {
	GetAll(ctx context.Context) ([]ucEntity.User, error)
}
