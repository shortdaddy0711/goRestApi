package main

import (
	"fmt"
	"net/http"

	"github.com/shortdaddy0711/go-rest-api/internal/comment"
	"github.com/shortdaddy0711/go-rest-api/internal/database"
	transportHTTP "github.com/shortdaddy0711/go-rest-api/internal/transport/http"
)

type App struct {}

func (app *App) Run() error {
	fmt.Println("Setting Up Our App")

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
		fmt.Println("Failed to setup server")
		return err
	}

	return nil
}

func main() {
	fmt.Println("Go REST API")
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting up our REST API")
		fmt.Println(err)
	}
}