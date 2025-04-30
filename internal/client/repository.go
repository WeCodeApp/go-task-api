package client

import (
	"context"
)

type ClientRepository interface {
	CreateClient(ctx context.Context, userID string, key []byte) error
	GetKeyById(ctx context.Context, userID string) ([]byte, error)
}
