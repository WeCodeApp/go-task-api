package service

import (
	"task/config"
	"task/internal/tasks"
)

const DEFAULT_LIMIT = 10

type TaskService struct {
	cfg      *config.Config
	taskRepo tasks.TaskRepository
}

func New(cfg *config.Config, taskRepo tasks.TaskRepository) *TaskService {
	return &TaskService{cfg: cfg, taskRepo: taskRepo}
}
