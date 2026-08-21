// Command gary is a tool-using coding agent that runs in the terminal.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"

	"github.com/orrybaram/gary/internal/agent"
	"github.com/orrybaram/gary/internal/tools"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "gary: %s\n", err)
		os.Exit(1)
	}
}

// run wires up dependencies and starts the agent loop. Keeping it separate
// from main means errors take one exit path and deferred cleanup still runs.
func run() error {
	client := anthropic.NewClient()

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	return agent.New(&client, getUserMessage, tools.All()).Run(context.Background())
}
