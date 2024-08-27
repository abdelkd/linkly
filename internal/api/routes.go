package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/abdelkd/linkly/internal/data"
)

type Application struct {
	Env struct {
		DatabaseURL string
		BaseURL     string
	}
	DB       *sql.DB
	Models   *data.Models
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Cfg      Config
}

type Config struct {
	Port int
}

func (app *Application) ParseEnv() error {
	envVars := []string{"DATABASE_URL", "BASE_URL"}
	for _, envString := range envVars {
		value := os.Getenv(envString)
		if value == "" {
			return fmt.Errorf("Environment variable %s is not set", envString)
		}

	}

	app.Env.BaseURL = os.Getenv("BASE_URL")
	app.Env.DatabaseURL = os.Getenv("DATABASE_URL")

	if !strings.HasSuffix(app.Env.BaseURL, "/") {
		app.Env.BaseURL = app.Env.BaseURL + "/"
	}

	return nil
}

func (app *Application) NewRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthcheck", app.handleHealthCheck)

	mux.HandleFunc("POST /v1/link", app.handleNewLink)
	mux.HandleFunc("/v1/link/{id}", app.handleGetLink)

	return mux
}

func (app *Application) NewServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Cfg.Port),
		Handler:      app.NewRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
