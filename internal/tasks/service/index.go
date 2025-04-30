package service

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"task/internal/models"
	httpErr "task/pkg/errors/http"
)

func (k *TaskService) Index(ctx *gin.Context) (resp models.IndexResponse, err error) {
	var filter models.ListFilter
	if l, ok := ctx.GetQuery("limit"); !ok {
		filter.Limit = DEFAULT_LIMIT
	} else {
		if i, err := strconv.Atoi(l); err != nil {
			return models.IndexResponse{}, httpErr.NewBadRequestError("Invalid limit")
		} else {
			filter.Limit = i
		}
	}
	if o, ok := ctx.GetQuery("offset"); ok {
		if i, err := strconv.Atoi(o); err != nil {
			return models.IndexResponse{}, httpErr.NewBadRequestError("Invalid offset")
		} else {
			filter.Offset = i
		}
	}

	count, err := k.taskRepo.Count(ctx)
	if err != nil {
		return models.IndexResponse{}, httpErr.NewInternalServerError("Unable to count tasks")
	}

	data, err := k.taskRepo.Get(ctx, filter)
	if err != nil {
		return models.IndexResponse{}, httpErr.NewInternalServerError("Unable to get tasks")
	}

	resp = models.IndexResponse{
		Data:  data,
		Total: count,
	}

	return resp, nil
}
