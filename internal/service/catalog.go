package service

import (
	"context"
	"log/slog"

	"base_app/internal/entity"
	"base_app/internal/usecase"
)

// CatalogService acts as a domain service for catalog operations.
type CatalogService struct {
	catalogRepo usecase.CatalogRepo
	log         *slog.Logger
}

// NewCatalogService creates a new CatalogService.
func NewCatalogService(repo usecase.CatalogRepo, log *slog.Logger) *CatalogService {
	return &CatalogService{
		catalogRepo: repo,
		log:         log,
	}
}

// GetCatalogItems retrieves all catalog items.
func (s *CatalogService) GetCatalogItems(ctx context.Context) ([]entity.CatalogItem, error) {
	// In a real app, more complex domain logic could go here.
	return s.catalogRepo.GetCatalogItems(ctx)
}
