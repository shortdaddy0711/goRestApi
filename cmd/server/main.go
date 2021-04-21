package main

import (
	"fmt"
	"net/http"

	"github.com/shortdaddy0711/go-rest-api/internal/comment"
	"github.com/shortdaddy0711/go-rest-api/internal/database"
	transportHTTP "github.com/shortdaddy0711/go-rest-api/internal/transport/http"

	log "github.com/sirupsen/logrus"
)

type App struct {
	Name    string
	Version string
}

func (app *App) Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		"AppName":    app.Name,
		"AppVersion": app.Version,
	}).Info("Setting up application")

	var err error
	db, err := database.NewDatabase()
	if err != nil {
		return err
	}
	err = database.MigrateDB(db)
	if err != nil {
		return err
	}

	commentService := comment.NewService(db)

	handler := transportHTTP.NewHandler(commentService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		log.Error("Failed to setup server")
		return err
	}

	log.Info("Now API up and running")
	return nil
}

func main() {
	fmt.Println("Go REST API")
	app := App{
		Name:    "Comment REST API",
		Version: "1.0.0",
	}
	if err := app.Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting up our REST API")
	}
}
