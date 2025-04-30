package client

import (
	"context"
	"gorm.io/gorm"
	"task/internal/core"
	"task/internal/models"
	"time"
)

type Repository struct {
	core.GormTransactionRepository
}

func New(db *gorm.DB) *Repository {
	return &Repository{core.NewRepository(db)}
}

func (c *Repository) CreateClient(ctx context.Context, userID string, key []byte) error {

	if err := c.Create(ctx, &models.Client{
		Issuer:    userID,
		Secret:    string(key),
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}

	return nil
}

func (c *Repository) GetKeyById(ctx context.Context, secret string) ([]byte, error) {
	var client models.Client
	if err := c.GetOneByField(ctx, &client, "issuer", secret); err != nil {
		return nil, err
	}

	return []byte(client.Secret), nil
}
