package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Environment             string `json:"environment"`
	Port                    string `json:"port"`
	DBBucketName            string `json:"db_bucket_name"`
	ImagesBucketName        string `json:"images_bucket_name"`
	FirebaseServiceAcctPath string `json:"firebase_service_account_path"`
}

// LoadConfig loads configuration from file or environment variables
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Try loading from config file first (optional, mainly for local dev)
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}

	// Try to load from file (don't error if it doesn't exist)
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables if present
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = env
	}
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}
	if bucket := os.Getenv("DB_BUCKET_NAME"); bucket != "" {
		config.DBBucketName = bucket
	}
	if bucket := os.Getenv("IMAGES_BUCKET_NAME"); bucket != "" {
		config.ImagesBucketName = bucket
	}
	if path := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH"); path != "" {
		config.FirebaseServiceAcctPath = path
	}

	// Set defaults
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.Environment == "" {
		config.Environment = "development"
	}

	// Validate required fields
	if config.DBBucketName == "" {
		return nil, fmt.Errorf("DB_BUCKET_NAME is required")
	}

	return config, nil
}
