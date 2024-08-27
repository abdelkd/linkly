package database

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/abdelkd/linkly/internal/api"
)

func InitDB(app *api.Application) {
	DB, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		app.ErrorLog.Fatalf("Failed to Create database connection: %v", err)
	}

	if err = DB.Ping(); err != nil {
		app.ErrorLog.Fatalf("Failed to connect to database: %v", err)
	}

	DB.SetConnMaxIdleTime(time.Hour)
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(100)

	app.DB = DB

	app.InfoLog.Println("Connected successfully to PostgreSQL")
}
