package usecase

import (
	"context"
	"log/slog"

	"base_app/internal/entity"
)

// CatalogUsecaseImpl handles the business logic for catalog operations.
type CatalogUsecaseImpl struct {
	service CatalogService
	log     *slog.Logger
}

// NewCatalogUsecase creates a new CatalogUsecase.
func NewCatalogUsecase(s CatalogService, l *slog.Logger) CatalogUsecase {
	return &CatalogUsecaseImpl{
		service: s,
		log:     l,
	}
}

// GetCatalogItems retrieves all catalog items.
func (uc *CatalogUsecaseImpl) GetCatalogItems(ctx context.Context) ([]entity.CatalogItem, error) {
	const op = "usecase.GetCatalogItems"

	items, err := uc.service.GetCatalogItems(ctx)
	if err != nil {
		uc.log.Error("failed to get catalog items", slog.String("op", op), slog.String("error", err.Error()))
		return nil, err
	}

	return items, nil
}
