package main

import (
	"fmt"
	"kv-db/internal/database"
	"os"

	"github.com/peterh/liner"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Starting")

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
	defer line.Close()

	line.SetCtrlCAborts(true)

	for {
		input, err := line.Prompt("> ")
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		if input == "exit" {
			break
		}

		if input != "" {
			line.AppendHistory(input)
		}

		result := db.HandleQuery(input)
		fmt.Println(result)
	}
}
