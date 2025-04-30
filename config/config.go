package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"log"
)

// App config struct
type Config struct {
	App    App
	Http   Http
	DB     DB
	Logger Logger
	Cors   Cors
}

type App struct {
	Version string `env:"VERSION,default=1.0.0"`
	Env     string `env:"ENV,default=local"`
	Debug   bool   `env:"DEBUG,default=true"`
	Log     string `env:"LOG,default=./tmp/api.log"`
}

type Http struct {
	Port string `env:"HTTP_PORT,default=:80"`
}

type DB struct {
	Host     string `env:"DB_HOST,default=localhost"`
	Port     string `env:"DB_PORT,default=3306"`
	DB       string `env:"DB_NAME,default=tasks"`
	User     string `env:"DB_USER,default=root"`
	Password string `env:"DB_PASSWORD,default=root"`
}

type Logger struct {
	DisableCaller     bool   `env:"LOGGER_DISABLE_CALLER,default=false"`
	DisableStacktrace bool   `env:"LOGGER_DISABLE_STACK_TRACE,default=false"`
	Encoding          string `env:"LOGGER_ENCODING,default=json"`
	Level             string `env:"LOGGER_LEVEL,default=info"`
}

type Cors struct {
	Whitelist string `env:"CORS_WHITELIST,default=http://localhost:8080"`
}

// Get config
func GetConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	var c Config
	if err := envdecode.Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
