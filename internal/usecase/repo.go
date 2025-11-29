package usecase

import (
	"base_app/internal/entity"
	"context"
)

// UserRepo is the interface for user database operations.
type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	// In a real app, you'd have more methods like CreateUser, etc.
}

// DataRepo is the interface for data database operations.
type DataRepo interface {
	SaveData(ctx context.Context, data *entity.Data) error
}

// CatalogRepo is the interface for catalog database operations.
type CatalogRepo interface {
	GetCatalogItems(ctx context.Context) ([]entity.CatalogItem, error)
}
