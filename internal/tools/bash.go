package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Bash runs a shell command and returns its combined output and exit code.
var Bash = Definition{
	Name: "bash",
	Description: `Run a shell command with bash and return its output.

Returns combined stdout and stderr, plus the exit code. Commands run in the
agent's working directory with the agent's own permissions. There is a default
timeout of 60 seconds; pass "timeout_seconds" to change it.`,
	InputSchema: schemaFor[bashInput](),
	Function:    runBash,
}

type bashInput struct {
	Command        string `json:"command" jsonschema_description:"The bash command to run."`
	TimeoutSeconds int    `json:"timeout_seconds,omitempty" jsonschema_description:"Optional timeout in seconds. Defaults to 60."`
}

const (
	// maxBashOutput caps how much output is fed back to the model.
	maxBashOutput = 30000
	// defaultBashTimeout applies when the caller doesn't set one.
	defaultBashTimeout = 60 * time.Second
)

func runBash(input json.RawMessage) (string, error) {
	var in bashInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}

	if strings.TrimSpace(in.Command) == "" {
		return "", fmt.Errorf("command must not be empty")
	}

	timeout := defaultBashTimeout
	if in.TimeoutSeconds > 0 {
		timeout = time.Duration(in.TimeoutSeconds) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", in.Command)
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
