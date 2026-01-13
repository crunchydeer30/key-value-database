package main

import (
	"fmt"
	"os"

	"github.com/crunchydeer30/key-value-database/internal/config"
	"github.com/crunchydeer30/key-value-database/internal/database"
	"github.com/crunchydeer30/key-value-database/internal/logger"
	"github.com/crunchydeer30/key-value-database/internal/network"
	"go.uber.org/zap"
)

var ConfigFileName = os.Getenv("CONFIG_FILE_NAME")

const DEFAULT_CONFIG_FILE_NAME = "config.yml"

func main() {
	if ConfigFileName == "" {
		ConfigFileName = DEFAULT_CONFIG_FILE_NAME
	}

	cfg, err := config.Load(ConfigFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config file: %v\n", err)
		os.Exit(1)
	}

	logger, err := logger.NewLogger(&cfg.Logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewDatabase(logger)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}

	server, err := network.NewTCPServer(
		cfg.Network.Address,
		db.HandleQuery,
		network.WithLogger(logger),
		network.WithMaxConnections(cfg.Network.MaxConnections),
	)
	if err != nil {
		logger.Fatal("failed to initialize network server", zap.Error(err))
	}

	server.Serve()
}
