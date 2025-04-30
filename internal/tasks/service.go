package tasks

import (
	"github.com/gin-gonic/gin"
	"task/internal/models"
)

type TaskService interface {
	Create(ctx *gin.Context) (resp *models.Task, err error)

	Delete(ctx *gin.Context) (id int, err error)

	Index(ctx *gin.Context) (resp models.IndexResponse, err error)

	Update(ctx *gin.Context) (resp *models.Task, err error)
}
