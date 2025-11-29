package postgresql

import (
	"context"
	"log/slog"

	"base_app/internal/adapter/repository/postgresql/sqlc"
	"base_app/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo implements the use case repository interfaces using sqlc.
type Repo struct {
	*sqlc.Queries
	pool *pgxpool.Pool
	log  *slog.Logger
}

// NewRepo creates a new repository.
func NewRepo(pool *pgxpool.Pool, log *slog.Logger) *Repo {
	return &Repo{
		Queries: sqlc.New(pool),
		pool:    pool,
		log:     log,
	}
}

// GetUserByEmail retrieves a user by their email address.
func (r *Repo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "adapter.sqlc.GetUserByEmail"

	userRow, err := r.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		r.log.Error("failed to get user by email", slog.String("op", op), slog.String("error", err.Error()))
		return nil, err
	}

	return &entity.User{
		ID:       userRow.ID,
		Email:    userRow.Email,
		Password: userRow.PasswordHash,
	}, nil
}

// SaveData saves data.
func (r *Repo) SaveData(ctx context.Context, data *entity.Data) error {
	const op = "adapter.sqlc.SaveData"

	err := r.Queries.SaveData(ctx, sqlc.SaveDataParams{
		Key:   data.Key,
		Value: data.Value,
	})
	if err != nil {
		r.log.Error("failed to save data", slog.String("op", op), slog.String("error", err.Error()))
		return err
	}
	return nil
}

// GetCatalogItems retrieves all catalog items.
func (r *Repo) GetCatalogItems(ctx context.Context) ([]entity.CatalogItem, error) {
	const op = "adapter.sqlc.GetCatalogItems"

	rows, err := r.Queries.GetCatalogItems(ctx)
	if err != nil {
		r.log.Error("failed to get catalog items", slog.String("op", op), slog.String("error", err.Error()))
		return nil, err
	}

	items := make([]entity.CatalogItem, len(rows))
	for i, row := range rows {
		items[i] = entity.CatalogItem{
			ID:          row.ID,
			Title:       row.Title,
			Description: row.Description.String,
			Disabled:    row.Disabled,
		}
	}

	return items, nil
}
