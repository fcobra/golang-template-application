package usecase

import (
	"context"
	"errors"
	"log/slog"

	"base_app/internal/entity"
	"base_app/pkg/hash"
)

// AuthUsecaseImpl handles the business logic for authentication.
type AuthUsecaseImpl struct {
	service AuthService
	log     *slog.Logger
}

// NewAuthUsecase creates a new AuthUsecase.
func NewAuthUsecase(s AuthService, l *slog.Logger) AuthUsecase {
	return &AuthUsecaseImpl{
		service: s,
		log:     l,
	}
}

// Authenticate finds a user by email and verifies their password.
func (uc *AuthUsecaseImpl) Authenticate(ctx context.Context, email, password string) (*entity.User, error) {
	const op = "usecase.Authenticate"

	user, err := uc.service.GetUserByEmail(ctx, email)
	if err != nil {
		uc.log.Error("failed to get user by email", slog.String("op", op), slog.String("error", err.Error()))
		return nil, err
	}

	if !hash.CheckPasswordHash(password, user.Password) {
		uc.log.Warn("invalid password attempt", slog.String("op", op), slog.String("email", email))
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
