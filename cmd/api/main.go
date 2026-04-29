package main

import (
	"log"

	"cryplio/internal/application"
	_ "github.com/lib/pq"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}
	defer app.DB.Close()

	if err := app.Router.Run(":" + app.Config.ServerPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
