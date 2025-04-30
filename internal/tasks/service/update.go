package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"task/internal/models"
	httpErr "task/pkg/errors/http"
)

func (k *TaskService) Update(ctx *gin.Context) (resp *models.Task, err error) {
	id := ctx.Param("id")

	var taskRequest models.TaskRequest
	if err := ctx.BindJSON(&taskRequest); err != nil {
		return nil, httpErr.NewBadRequestError("Request data not properly formatted")
	}

	taskFilter := map[string]interface{}{
		"id =":      id,
		"deleted =": false,
	}

	existingTask, err := k.taskRepo.GetById(ctx, taskFilter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpErr.NewNotFoundError("Task not found")
		}
		return nil, httpErr.NewInternalServerError("Unable to get task")
	}

	var taskData models.Task
	err = copier.Copy(&taskData, &taskRequest)
	if err != nil {
		return nil, httpErr.NewInternalServerError("Unable to copy data")
	}

	taskData.ID = existingTask.ID

	err = k.taskRepo.Update(ctx, &taskData)
	if err != nil {
		return nil, httpErr.NewInternalServerError("Unable to update")
	}

	return &taskData, nil
}
