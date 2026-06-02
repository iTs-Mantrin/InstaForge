package main

import (
	"fmt"
	"log"
	"os"

	"instaforge/internal/config"
	"instaforge/internal/database"
	"instaforge/internal/logger"
)

func main() {
	cfg := config.Load()

	logger.Init(cfg.App)

	db, err := database.RetryableConnect(cfg.DB, 3)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fmt.Println("Migrations completed successfully")
	os.Exit(0)
}
