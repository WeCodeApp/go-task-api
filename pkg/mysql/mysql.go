package mysql

import (
	gosql "github.com/go-sql-driver/mysql"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	mysqld "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"task/config"
	"time"
)

func New(conf *config.Config) (*gorm.DB, error) {
	cfg := &gosql.Config{
		Net:                  "tcp",
		Addr:                 conf.DB.Host,
		DBName:               conf.DB.DB,
		User:                 conf.DB.User,
		Passwd:               conf.DB.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	logLevel := logger.Error
	if conf.App.Debug {
		logLevel = logger.Info
	}
	db, err := gorm.Open(mysqld.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, err
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(time.Minute * 3)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
