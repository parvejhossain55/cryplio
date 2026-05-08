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

	// Start Background Worker
	// go func() {
	// 	log.Printf("starting background worker...")
	// 	if err := app.Worker.Start(); err != nil {
	// 		log.Printf("worker failed: %v", err)
	// 	}
	// }()

	// Start Task Scheduler
	// go func() {
	// 	log.Printf("starting task scheduler...")
	// 	if err := app.Scheduler.Start(); err != nil {
	// 		log.Printf("scheduler failed: %v", err)
	// 	}
	// }()

	log.Printf("HTTP server on port %s", app.Config.ServerPort)
	if err := app.Router.Run(":" + app.Config.ServerPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
