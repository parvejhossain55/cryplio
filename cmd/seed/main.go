package main

import (
	"context"
	"log"

	"cryplio/pkg/config"
	"cryplio/pkg/database"
	"cryplio/seeder"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Open(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Initialize and run the seeder
	s := seeder.NewSeeder(db, cfg)
	if err := s.SeedAll(ctx); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}
}
