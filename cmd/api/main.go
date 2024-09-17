package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/abdelkd/linkly/internal/api"
	"github.com/abdelkd/linkly/internal/data"
	"github.com/abdelkd/linkly/internal/database"
)

func main() {
	var cfg api.Config

	flag.IntVar(&cfg.Port, "port", 4000, "The default port for the server API")
	flag.Parse()
	fmt.Println(cfg)

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ltime|log.Ldate)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate)

	application := &api.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		Cfg:      cfg,
	}

	database.InitDB(application)

	application.Models = data.NewModels(application.DB)
	server := application.NewServer()

	err := application.ParseEnv()
	if err != nil {
		errorLog.Fatal(err)
	}

	infoLog.Printf("Listening on port :%d\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}
