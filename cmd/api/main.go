package main

import (
	"os"

	"cryplio/internal/application"
	"cryplio/pkg/logger"

	_ "github.com/lib/pq"
)

func main() {
	app, err := application.New()
	if err != nil {
		logger.Error("bootstrap failed", logger.Fields{"error": err.Error()})
		os.Exit(1)
	}
	defer app.DB.Close()

	// Start Background Worker
	// go func() {
	// 	logger.Info("starting background worker", logger.Fields{})
	// 	if err := app.Worker.Start(); err != nil {
	// 		logger.Error("worker failed", logger.Fields{"error": err.Error()})
	// 	}
	// }()

	// Start Task Scheduler
	// go func() {
	// 	logger.Info("starting task scheduler", logger.Fields{})
	// 	if err := app.Scheduler.Start(); err != nil {
	// 		logger.Error("scheduler failed", logger.Fields{"error": err.Error()})
	// 	}
	// }()

	logger.Info("starting HTTP server", logger.Fields{"port": app.Config.ServerPort})
	if err := app.Router.Run(":" + app.Config.ServerPort); err != nil {
		logger.Error("server failed", logger.Fields{"error": err.Error()})
		os.Exit(1)
	}
}
