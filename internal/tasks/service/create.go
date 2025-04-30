package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"task/internal/models"
	httpErr "task/pkg/errors/http"
)

func (k *TaskService) Create(ctx *gin.Context) (resp *models.Task, err error) {
	var t models.TaskRequest
	if err := ctx.BindJSON(&t); err != nil {
		return nil, httpErr.NewBadRequestError("Request data not properly formatted")
	}

	var b models.Task
	err = copier.Copy(&b, &t)
	if err != nil {
		return nil, httpErr.NewInternalServerError("Unable to copy data")
	}

	newTask, err := k.taskRepo.Create(ctx, &b)
	if err != nil {
		return nil, httpErr.NewInternalServerError("Unable to create data")
	}

	return newTask, nil
}
