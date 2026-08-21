package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}

		return scanner.Text(), true
	}
	tools := []ToolDefinition{ReadFileDefinition, ListFilesDefinition, EditFileDefinition, BashDefinition}
	agent := NewAgent(&client, getUserMessage, tools)
	err := agent.Run(context.TODO())

	if err != nil {
		fmt.Printf("Err %s\n", err.Error())
	}
}
