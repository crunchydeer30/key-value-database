//nolint:forbidigo
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/crunchydeer30/key-value-database/internal/network"
	"github.com/peterh/liner"
)

func main() {
	address := flag.String("address", "localhost:3223", "address of the server")
	flag.Parse()

	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	//nolint:errcheck
	defer line.Close()

	fmt.Println("Connecting to server at", *address, "...")

	client, err := network.NewTCPClient(*address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to server: %v\n", err)
		return
	}
	//nolint:errcheck
	defer client.Close()

	fmt.Println("Connected to server")

	for {
		input, err := line.Prompt("> ")
		if err != nil {
			if !errors.Is(err, liner.ErrPromptAborted) {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			}

			break
		}

		if input == "exit" {
			break
		}

		if input == "clear" {
			_, err := os.Stdout.WriteString("\033[H\033[2J")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error clearing screen: %v\n", err)
				continue
			}
			continue
		}

		if input != "" {
			line.AppendHistory(input)
		}

		result, err := client.Send([]byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			continue
		}

		//nolint:forbidigo
		fmt.Print(string(result))
	}
}
