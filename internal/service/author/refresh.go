package author

import (
	"context"
	"errors"
	ucEntity "go-projects/hexagonal-example/internal/service/entity/author"
)

// Refresh menukar refresh token dengan pasangan token baru (stateless: refresh
// token tidak disimpan, cukup diverifikasi tanda tangan & masa berlakunya).
func (s service) Refresh(ctx context.Context, req ucEntity.RefreshRequest) (ucEntity.AuthResponse, error) {
	var response ucEntity.AuthResponse

	claims, err := s.token.ParseRefresh(req.RefreshToken)
	if err != nil {
		return response, errors.New("refresh token tidak valid atau kedaluwarsa")
	}

	access, refresh, expiresIn, err := s.token.GeneratePair(claims.UserID, claims.Role)
	if err != nil {
		return response, err
	}

	return ucEntity.AuthResponse{
		TokenType:    "Bearer",
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		Role:         claims.Role,
	}, nil
}

type IRefresh interface {
	Refresh(ctx context.Context, req ucEntity.RefreshRequest) (ucEntity.AuthResponse, error)
}
