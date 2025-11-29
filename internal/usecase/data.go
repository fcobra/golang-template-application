package usecase

import (
	"context"
	"errors"
	"log/slog"

	"base_app/internal/entity"
)

// DataUsecaseImpl handles the business logic for data operations.
type DataUsecaseImpl struct {
	service DataService
	log     *slog.Logger
}

// NewDataUsecase creates a new DataUsecase.
func NewDataUsecase(s DataService, l *slog.Logger) DataUsecase {
	return &DataUsecaseImpl{
		service: s,
		log:     l,
	}
}

// SaveData validates and saves data.
func (uc *DataUsecaseImpl) SaveData(ctx context.Context, data *entity.Data) error {
	const op = "usecase.SaveData"

	if data.Key == "" {
		return errors.New("key cannot be empty")
	}

	err := uc.service.SaveData(ctx, data)
	if err != nil {
		uc.log.Error("failed to save data", slog.String("op", op), slog.String("error", err.Error()))
		return err
	}

	uc.log.Info("data saved successfully", slog.String("op", op), slog.String("key", data.Key))
	return nil
}
