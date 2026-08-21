# Gary

A small, tool-using coding agent that lives in your terminal.

Gary talks to the Anthropic API, and when the model asks for a tool, Gary runs
it locally: reading files, listing directories, editing text, and executing
bash commands in your working directory.

## Requirements

- Go 1.27+
- An Anthropic API key

## Setup

    export ANTHROPIC_API_KEY=sk-ant-...

## Build and run

    make build     # builds bin/gary
    make run       # builds, then starts the agent
    make test      # go test ./...
    make fmt       # gofmt -w .
    make vet       # go vet ./...
    make tidy      # go mod tidy
    make clean     # remove bin/

Or directly:

    go run ./cmd/gary

## Usage

Start Gary and type at the `You:` prompt. He replies inline, and any tool call
is echoed as `tool: name({...})` before it runs. Ctrl-C quits.

    You: what does the spinner do when stdout is piped?
    tool: read_file({"path":"internal/ui/spinner.go"})
    Gary: It returns an inert spinner, so escape codes never end up in
    machine-readable output.

## Tools

| Tool | What it does |
| --- | --- |
| `read_file` | Read the contents of a relative file path. |
| `list_files` | Recursively list entries under a path (defaults to `.`). |
| `edit_file` | Replace an exact substring in a file; creates the file when `old_str` is empty. |
| `bash` | Run a shell command, returning combined output and exit code. 60s default timeout, output capped at 30k chars. |

## Layout

    cmd/gary/main.go      entrypoint; wires up client, stdin reader, toolset
    internal/agent/       conversation loop, tool dispatch, system prompt
    internal/tools/       tool definitions, JSON schemas, registry
    internal/ui/          banner, colors, spinner

## Adding a tool

1. Add a file in `internal/tools/` defining a `Definition` (name, description,
   input struct with `jsonschema_description` tags, and the Go function).
2. Append it to the slice returned by `All()` in `internal/tools/registry.go`.

The input schema is reflected from your struct via `schemaFor[T]()`, so the
struct is the single source of truth for what the model sees.

## A word of caution

Gary executes shell commands and writes files with your permissions and no
confirmation prompt. Run him in a directory you're happy to have modified,
ideally one under version control.
