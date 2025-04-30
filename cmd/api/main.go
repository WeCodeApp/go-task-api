package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"task/config"
	_ "task/docs"
	"task/internal/server"
	"task/pkg/logger"
	"task/pkg/mysql"
)

// @title         Task API
// @version       1.0.0
// @description   Simple Task Manager API
// @contact.name  shairayvonnecruz@gmail.com
// @BasePath      /api
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	var file io.Writer
	file = os.Stderr
	if cfg.App.Debug && cfg.App.Log != "" {
		filename := cfg.App.Log
		if filename == "" {
			filename = "./tmp/api.log"
		}
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Opening file: %v", err)
		}
	}

	l := logger.New(cfg, file)
	l.Infof(
		"AppVersion: %s, LogLevel: %s, Environment: %s",
		cfg.App.Version,
		cfg.Logger.Level,
		cfg.App.Env,
	)

	app := newApp(cfg, l)
	if err := app.setup(); err != nil {
		l.Errorf("could not setup app: %s", err)
		return
	}
	defer app.close()

	if err := app.run(); err != nil {
		l.Errorf("could not run app: %s", err)
	}
}

type Application struct {
	conf    *config.Config
	logger  logger.Logger
	closeFn func() error
	server  *server.Server
}

func newApp(c *config.Config, l logger.Logger) *Application {
	return &Application{conf: c, logger: l}
}

func (app *Application) setup() error {
	db, err := mysql.New(app.conf)
	if err != nil {
		return fmt.Errorf("cannot setup db: %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("cannot setup db: %s", err)
	}

	server, err := server.New(app.logger, app.conf, db)
	if err != nil {
		return fmt.Errorf("cannot setup server: %s", err)
	}
	app.server = server
	app.closeFn = func() (err error) {
		return sqlDB.Close()
	}
	return nil
}

func (app *Application) run() error {
	return app.server.Run()
}

func (app *Application) close() {
	if err := app.closeFn(); err != nil {
		app.logger.Error("error when closing app: %s", err)
	}
}
