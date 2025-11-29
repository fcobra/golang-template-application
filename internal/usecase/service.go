package usecase

import (
	"base_app/internal/entity"
	"context"
)

// AuthService defines the interface for the authentication domain service.
type AuthService interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

// DataService defines the interface for the data domain service.
type DataService interface {
	SaveData(ctx context.Context, data *entity.Data) error
}

// CatalogService defines the interface for the catalog domain service.
type CatalogService interface {
	GetCatalogItems(ctx context.Context) ([]entity.CatalogItem, error)
}
