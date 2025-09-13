package main

import (
	"log"
	"uptime-monitor/internal/api"
	"uptime-monitor/internal/config"
	"uptime-monitor/internal/database"
	"uptime-monitor/internal/monitoring"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// 2. Initialize database connection
	dbPool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("FATAL: could not connect to database: %v", err)
	}
	defer dbPool.Close()
	log.Println("INFO: Database connection successful")

	// 3. Start the monitoring worker in the background
	log.Println("INFO: Initializing monitoring worker...")
	monitor := monitoring.NewMonitor(dbPool)
	go monitor.Start() // Starts the worker in a new goroutine

	// 4. Start the API server (this is a blocking call)
	server := api.NewServer(dbPool)
	log.Printf("INFO: Starting API server on %s", cfg.ServerAddress)
	if err := server.Start(cfg.ServerAddress); err != nil {
		log.Fatalf("FATAL: could not start server: %v", err)
	}
}
