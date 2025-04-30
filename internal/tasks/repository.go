package tasks

import (
	"context"
	"task/internal/models"
)

type TaskRepository interface {
	Count(ctx context.Context) (count int64, err error)

	Create(ctx context.Context, task *models.Task) (*models.Task, error)

	Delete(ctx context.Context, task *models.Task) error

	Get(ctx context.Context, filter models.ListFilter) (tasks []*models.Task, err error)

	GetById(ctx context.Context, filters map[string]interface{}) (task *models.Task, err error)

	Update(ctx context.Context, task *models.Task) error
}
