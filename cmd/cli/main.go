package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/crunchydeer30/key-value-database/internal/config"
	"github.com/crunchydeer30/key-value-database/internal/database"
	"github.com/crunchydeer30/key-value-database/internal/logger"
	"github.com/peterh/liner"
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

	logger, err := logger.NewLogger(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewDatabase(logger)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}

	line := liner.NewLiner()
	defer func() {
		if err := line.Close(); err != nil {
			logger.Error("failed to close liner", zap.Error(err))
		}
	}()

	line.SetCtrlCAborts(true)

	for {
		input, err := line.Prompt("> ")
		if err != nil {
			if !errors.Is(err, liner.ErrPromptAborted) {
				logger.Error("Error reading input:", zap.Error(err))
			}

			break
		}

		if input == "exit" {
			break
		}

		if input != "" {
			line.AppendHistory(input)
		}

		result := db.HandleQuery(input)
		//nolint:forbidigo
		fmt.Println(result)
	}
}
