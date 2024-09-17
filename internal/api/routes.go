package api

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/abdelkd/linkly/internal/data"
	"github.com/gorilla/mux"
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

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s - %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (app *Application) NewRoutes() http.Handler {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/v1/").Subrouter()
	apiRouter.HandleFunc("/healthcheck", app.handleHealthCheck)

	apiRouter.HandleFunc("/link", app.handleNewLink)
	apiRouter.HandleFunc("/link/{id}", app.handleGetLink)

	// Handle Static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.HandleFunc("/{id}", app.handleGetLink)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = templ.Execute(w, nil)
		if err != nil {
			app.serverError(w, err)
			return
		}
	})

	return loggerMiddleware(router)
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
