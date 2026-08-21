package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var BashDefinition = ToolDefinition{
	Name: "bash",
	Description: `Run a shell command with bash and return its output.

Returns combined stdout and stderr, plus the exit code. Commands run in the
agent's working directory with the agent's own permissions. There is a default
timeout of 60 seconds; pass "timeout_seconds" to change it.`,
	InputSchema: BashInputSchema,
	Function:    Bash,
}

type BashInput struct {
	Command        string `json:"command" jsonschema_description:"The bash command to run."`
	TimeoutSeconds int    `json:"timeout_seconds,omitempty" jsonschema_description:"Optional timeout in seconds. Defaults to 60."`
}

var BashInputSchema = GenerateSchema[BashInput]()

const maxBashOutput = 30000

func Bash(input json.RawMessage) (string, error) {
	bashInput := BashInput{}
	if err := json.Unmarshal(input, &bashInput); err != nil {
		return "", err
	}

	if strings.TrimSpace(bashInput.Command) == "" {
		return "", fmt.Errorf("command must not be empty")
	}

	timeout := 60 * time.Second
	if bashInput.TimeoutSeconds > 0 {
		timeout = time.Duration(bashInput.TimeoutSeconds) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", bashInput.Command)
	output, err := cmd.CombinedOutput()

	result := string(output)
	if len(result) > maxBashOutput {
		result = result[:maxBashOutput] + "\n... [output truncated]"
	}

	if ctx.Err() == context.DeadlineExceeded {
		return result + fmt.Sprintf("\n[command timed out after %s]", timeout), nil
	}

	if cmd.ProcessState == nil {
		return "", fmt.Errorf("failed to run command: %w", err)
	}

	return fmt.Sprintf("exit code: %d\n%s", cmd.ProcessState.ExitCode(), result), nil
}
