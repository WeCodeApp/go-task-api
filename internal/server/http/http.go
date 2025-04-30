package http

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"
	"strings"
	"task/config"
	"task/internal/middlewares"
	"task/internal/tasks"
	handler "task/internal/tasks/delivery/http"
	httpErrors "task/pkg/errors/http"
	"task/pkg/gin-contrib/recover"
	"task/pkg/gin-contrib/secure"
	"time"
)

const (
	shutdownTimeout = 10 * time.Second
	maxHeaderBytes  = 1 << 20
	gzipLevel       = 5
	stackSize       = 1 << 10
	readTimeout     = 5 * time.Second
	writeTimeout    = 5 * time.Second
	local           = "local"
)

type Server struct {
	cfg    *config.Config
	server *http.Server
}

func New(cfg *config.Config, m *middlewares.Middleware, taskService tasks.TaskService) (*Server, error) {
	router := gin.New()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, httpErrors.New(http.StatusNotFound, "Route Not Found", nil))
	})
	router.Use(requestid.New())
	router.Use(gzip.Gzip(gzipLevel, gzip.WithExcludedPaths([]string{"/swagger/"})))
	router.Use(secure.Secure())
	if strings.ToLower(cfg.App.Env) == local {
		router.Use(cors.New(CORSConfig(cfg)))
	}

	router.Use(recover.RecoverWithConfig(recover.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	api := router.Group("/api", m.Authorize)
	handler := handler.New(cfg, taskService)
	handler.Setup(api)
	router.GET("/swagger/*any", gs.WrapHandler(swaggerfiles.Handler))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	return &Server{
		cfg,
		&http.Server{
			Addr:           cfg.Http.Port,
			Handler:        router,
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			MaxHeaderBytes: maxHeaderBytes,
		},
	}, nil
}

func CORSConfig(cfg *config.Config) cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = strings.Split(cfg.Cors.Whitelist, ",")
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization", "X-API-Key")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}

// Run starts serving and listening http server with graceful shutdown.
func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
