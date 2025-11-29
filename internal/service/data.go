package service

import (
	"context"
	"base_app/internal/entity"
	"base_app/internal/usecase"
	"log/slog"
)

// DataService acts as a domain service for data operations.
type DataService struct {
	dataRepo usecase.DataRepo
	log      *slog.Logger
}

func NewDataService(dataRepo usecase.DataRepo, log *slog.Logger) *DataService {
	return &DataService{
		dataRepo: dataRepo,
		log:      log,
	}
}

func (s *DataService) SaveData(ctx context.Context, data *entity.Data) error {
	// In a real app, more complex domain logic could go here.
	return s.dataRepo.SaveData(ctx, data)
}
