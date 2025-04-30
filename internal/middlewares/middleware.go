package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task/config"
	httpErr "task/pkg/errors/http"
	"task/pkg/guard"
)

type Middleware struct {
	cfg   *config.Config
	guard *guard.Guard
}

func New(cfg *config.Config, g *guard.Guard) *Middleware {
	return &Middleware{cfg: cfg, guard: g}
}

func (m *Middleware) Authorize(c *gin.Context) {
	key := c.Request.Header.Get("X-API-Key")

	userId, err := m.guard.IsKeyValid(c, []byte(key))
	if err {
		c.Set("Client", userId)
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, httpErr.NewUnauthorizedError("invalid key"))
		c.Abort()
	}
}
