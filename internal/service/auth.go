package service

import (
	"context"
	"base_app/internal/entity"
	"base_app/internal/usecase" // To get the repo interface
	"log/slog"
)

// AuthService acts as a domain service for authentication.
// In this simple case, it's a pass-through to the repository.
// In a more complex app, it could contain logic shared between different auth-related usecases.
type AuthService struct {
	userRepo usecase.UserRepo
	log      *slog.Logger
}

func NewAuthService(userRepo usecase.UserRepo, log *slog.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		log:      log,
	}
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}
