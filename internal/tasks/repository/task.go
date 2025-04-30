package repository

import (
	"context"
	"gorm.io/gorm"
	"task/internal/core"
	"task/internal/models"
)

type TaskRepository struct {
	db core.GormTransactionRepository
}

func New(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: core.NewRepository(db),
	}
}

func (r *TaskRepository) Count(ctx context.Context) (count int64, err error) {
	var taskModels []*models.Task
	err = r.db.GetWhereCount(ctx, &taskModels, &count, "deleted = false")
	return
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	if err := r.db.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (r *TaskRepository) Delete(ctx context.Context, task *models.Task) error {
	return r.db.Save(ctx, task)
}

func (r *TaskRepository) Get(ctx context.Context, filter models.ListFilter) (tasks []*models.Task, err error) {
	err = r.db.GetWhereBatch(ctx, &tasks, "deleted = false", filter.Limit, filter.Offset)
	return
}

func (r *TaskRepository) GetById(ctx context.Context, filters map[string]interface{}) (task *models.Task, err error) {
	err = r.db.GetOneByFields(ctx, &task, filters)
	return
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.Save(ctx, task)
}
