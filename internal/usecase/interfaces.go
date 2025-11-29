package usecase

import (
	"base_app/internal/entity"
	"context"
)

// AuthUsecase defines the interface for authentication business logic.
type AuthUsecase interface {
	Authenticate(ctx context.Context, email, password string) (*entity.User, error)
}

// DataUsecase defines the interface for data-related business logic.
type DataUsecase interface {
	SaveData(ctx context.Context, data *entity.Data) error
}

// CatalogUsecase defines the interface for catalog-related business logic.
type CatalogUsecase interface {
	GetCatalogItems(ctx context.Context) ([]entity.CatalogItem, error)
}
