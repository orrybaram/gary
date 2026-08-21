package ui

import (
	"fmt"
	"os"
)

// ANSI escape codes used across the terminal UI. These stay unexported: all
// styling happens through the helpers in this package so escape sequences
// never leak into business logic.
const (
	colorReset = "\u001b[0m"
	colorBold  = "\u001b[1m"
	colorDim   = "\u001b[2m"

	colorUser  = "\u001b[94m"       // blue
	colorWarn  = "\u001b[93m"       // yellow
	colorAgent = "\u001b[38;5;208m" // orange
	colorTool  = "\u001b[92m"       // green
	colorCat   = "\u001b[38;5;208m" // orange
)

// UserPrompt writes the "You: " prompt the user types against.
func UserPrompt() {
	fmt.Printf("%sYou%s: ", colorUser, colorReset)
}

// AgentMessage prints a line of assistant prose, labelled with the agent name.
func AgentMessage(name, text string) {
	fmt.Printf("%s%s%s: %s\n", colorAgent, name, colorReset, text)
}

// Warn prints a diagnostic about the agent loop itself, not about the code
// being worked on. Goes to stderr so it stays out of piped transcripts.
func Warn(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%swarning%s: %s\n", colorWarn, colorReset, fmt.Sprintf(format, args...))
}

// ToolCall prints the name and raw JSON input of a tool the model invoked.
func ToolCall(name, input string) {
	fmt.Printf("%stool%s: %s(%s)\n", colorTool, colorReset, name, input)
}
