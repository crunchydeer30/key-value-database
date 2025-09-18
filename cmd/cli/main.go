package main

import (
	"fmt"
	"kv-db/internal/database"
	"log"

	"go.uber.org/zap"
)

func main() {
	fmt.Println("Starting")

	logger, err := zap.NewDevelopment()

	if err != nil {
		log.Fatalf("failed to initialize logger %v", err)
	}

	db, err := database.NewDatabase(logger)
	if err != nil {
		log.Fatalf("failed to initialize db %v", err)
	}

	fmt.Println("Database initialized")

	err = db.HandleQuery("GET arg1 arg2")

	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Print("query handled")
}
