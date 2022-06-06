package main

import (
	"log"
	"os"

	"logging/internal/config"
	"logging/internal/randlogs"
	"logging/pkg/logging"
)

func main() {
	cfg, err := config.Load("./etc/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %+v", err)
		return
	}

	credentialsJSON, err := os.ReadFile(cfg.GCL.ServiceAccountPath)
	if err != nil {
		log.Fatalf("failed to read file %s: %+v", cfg.GCL.ServiceAccountPath, err)
		return
	}

	logger, err := logging.New(logging.Options{
		Level: cfg.Logging.Level,
		GCL: logging.GoogleCloudLoggingOptions{
			CredentialsJSON: credentialsJSON,
			ProjectID:       cfg.GCL.ProjectID,
			LogID:           cfg.GCL.LogID,
			Level:           cfg.GCL.Level,
		},
	})
	if err != nil {
		log.Fatalf("failed to build logger: %+v", err)
		return
	}

	randlogs.StartLogging(logger)
}
