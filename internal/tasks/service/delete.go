package service

import (
	"github.com/gin-gonic/gin"
	"strconv"
	httpErr "task/pkg/errors/http"
)

func (k *TaskService) Delete(ctx *gin.Context) (id int, err error) {
	idStr := ctx.Param("id")
	id, err = strconv.Atoi(idStr)
	if err != nil {
		return 0, httpErr.NewBadRequestError("Invalid ID")
	}

	taskFilter := map[string]interface{}{
		"id =":      id,
		"deleted =": false,
	}

	existingTask, err := k.taskRepo.GetById(ctx, taskFilter)
	if err != nil {
		return 0, httpErr.NewNotFoundError("No tasks found with this id")
	}

	existingTask.Deleted = true

	err = k.taskRepo.Delete(ctx, existingTask)
	if err != nil {
		return 0, httpErr.NewInternalServerError("Unable to delete")
	}

	return id, nil
}
