package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pablo-puyat/elinor/internal/api"
	"github.com/pablo-puyat/elinor/internal/config"
	"github.com/pablo-puyat/elinor/internal/monitor"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := config.InitLogger(cfg.LogFile)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize monitor
	networkMonitor := monitor.New(logger)

	// Start the monitor
	go networkMonitor.Start(cfg.UpdateInterval)

	// Initialize and start API server
	apiServer := api.New(cfg.APIPort, networkMonitor, logger)
	go apiServer.Start()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-sigChan
	logger.Printf("Received signal %v, shutting down...", sig)

	// Cleanup
	apiServer.Stop()
	networkMonitor.Stop()
}
