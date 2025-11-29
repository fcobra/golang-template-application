package inmemory

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"base_app/internal/entity"
	"base_app/pkg/hash"

	"github.com/google/uuid"
)

// Adapter implements the AuthService interface with an in-memory store.
type Adapter struct {
	log *slog.Logger
}

// New creates a new in-memory auth adapter.
func New(log *slog.Logger) *Adapter {
	return &Adapter{
		log: log,
	}
}

var hardcodedUser = entity.User{
	ID:        uuid.New(),
	Email:     "test@example.com",
	Password:  "$2a$10$WhWf0qQzwtD8fz6p/Ge.2e8Y6WhZRN/vopNJXofJ7vEaG4KEukRPS", // "password123"
	CreatedAt: time.Now(),
}

// GetUserByEmail simulates fetching a user from an in-memory store.
func (a *Adapter) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "adapter.inmemory.GetUserByEmail"

	if email == hardcodedUser.Email {
		a.log.Info("found user in in-memory store", slog.String("op", op), slog.String("email", email))
		return &hardcodedUser, nil
	}

	a.log.Warn("user not found in in-memory store", slog.String("op", op), slog.String("email", email))
	return nil, errors.New("user not found")
}

// Authenticate is a dummy implementation for the in-memory adapter.
// The actual password check happens in the usecase.
func (a *Adapter) Authenticate(ctx context.Context, email, password string) (*entity.User, error) {
	if email == hardcodedUser.Email && hash.CheckPasswordHash(password, hardcodedUser.Password) {
		return &hardcodedUser, nil
	}
	return nil, errors.New("invalid credentials")
}
