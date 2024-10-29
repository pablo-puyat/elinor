package config

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	UpdateInterval time.Duration
	LogFile        string
	APIPort        int
}

func Load() (*Config, error) {
	// In a real application, you might want to load this from a file
	return &Config{
		UpdateInterval: 5 * time.Second,
		LogFile:        "/var/log/netmonitor.log",
		APIPort:        8080,
	}, nil
}

func InitLogger(logFile string) (*log.Logger, error) {
	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// Open log file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return log.New(file, "", log.LstdFlags), nil
}
