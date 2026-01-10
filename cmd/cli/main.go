package main

import (
	"fmt"
	"os"

	"github.com/crunchydeer30/key-value-database/internal/database"
	"github.com/peterh/liner"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
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
			logger.Error("Error reading input:", zap.Error(err))
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
