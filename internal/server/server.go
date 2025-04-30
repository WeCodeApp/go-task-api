package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	taskRepository "task/internal/tasks/repository"
	taskService "task/internal/tasks/service"
	"task/pkg/logger"

	"github.com/gin-gonic/gin"

	// "github.com/opentracing/opentracing-go"
	"task/internal/middlewares"
	"task/pkg/guard"

	"gorm.io/gorm"
	"task/config"
	"task/internal/client/repository"

	task "task/internal/server/http"
	//"time"
)

// Server
type Server struct {
	logger     logger.Logger
	cfg        *config.Config
	taskServer *task.Server
}

// Server constructor
func New(logger logger.Logger, cfg *config.Config, db *gorm.DB) (*Server, error) {
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	repository := taskRepository.New(db)
	service := taskService.New(cfg, repository)

	clientRepository := client.New(db)
	g, err := guard.New(
		clientRepository,
		guard.WithKeyLen(guard.DefaultKeyLength),
	)

	if err != nil {
		return nil, err
	}

	taskServer, err := task.New(cfg, middlewares.New(cfg, g), service)
	if err != nil {
		return nil, err
	}
	return &Server{logger, cfg, taskServer}, nil
}

// Run server
func (s *Server) Run() error {
	s.logger.Info("Server Starting...")
	return s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		s.taskServer.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.logger.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.logger.Errorf("ctx.Done: %v", done)
	}

	if err := s.taskServer.Shutdown(ctx); err != nil {
		s.logger.Errorf("HTTP Server.Shutdown: %v", err)
	}

	s.logger.Info("Server Exited Properly")
	return nil
}
