package http

import (
	"task/internal/models"
	"task/internal/tasks"
	httpErr "task/pkg/errors/http"

	//"context"
	"github.com/gin-gonic/gin"
	"net/http"
	//"net/url"
	"task/config"
)

type CommonResponse struct {
	Message string `json:"message,omitempty"`
}

type Handler struct {
	taskService tasks.TaskService
}

func New(
	cfg *config.Config,
	taskService tasks.TaskService) *Handler {
	res := &Handler{
		taskService: taskService,
	}
	return res
}

func (h *Handler) Setup(g *gin.RouterGroup) {
	g.GET("/tasks", h.Index)
	g.POST("/tasks", h.Create)
	g.PATCH("/tasks/:id", h.Update)
	g.DELETE("/tasks/:id", h.Delete)
}

// Task godoc
// @Summary Get all saved tasks
// @Tags V1
// @Description Get tasks
// @Param limit query string false "limit"
// @Param offset query string false "offset"
// @Param X-API-Key header string true "API Key"
// @Success 200 {object}  models.IndexResponse
// @Router /tasks [get]
func (h *Handler) Index(c *gin.Context) {
	resp, err := h.taskService.Index(c)
	if err != nil {
		httpErr.HTTPErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

// Task godoc
// @Summary Create task in database
// @Tags V1
// @Description Create task
// @Accept json
// @Param data body models.TaskRequest true "Task Data"
// @Param X-API-Key header string true "API Key"
// @Success 200 {object}  models.TaskResponse
// @Router /tasks [post]
func (h *Handler) Create(c *gin.Context) {
	resp, err := h.taskService.Create(c)
	if err != nil {
		httpErr.HTTPErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

// Task godoc
// @Summary Update task in database
// @Tags V1
// @Description Update task
// @Accept json
// @Param data body models.TaskRequest true "Task Data"
// @Param X-API-Key header string true "API Key"
// @Param id path int true "Task ID"
// @Success 200 {object} models.TaskResponse
// @Router /tasks/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	resp, err := h.taskService.Update(c)
	if err != nil {
		httpErr.HTTPErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

// Task godoc
// @Summary Delete task in database
// @Tags V1
// @Description Delete task
// @Accept json
// @Param X-API-Key header string true "API Key"
// @Param id path int true "Task ID"
// @Success 200 {object} models.DeleteResponse
// @Router /tasks/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := h.taskService.Delete(c)
	if err != nil {
		httpErr.HTTPErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, models.DeleteResponse{ID: id, Status: "deleted"})
	return
}
