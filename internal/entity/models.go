package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // The password hash, ignored by json marshalling
	CreatedAt time.Time `json:"created_at"`
}

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CatalogItem struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Disabled    bool      `json:"disabled"`
}
