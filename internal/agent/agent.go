package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"

	"github.com/orrybaram/gary/internal/tools"
	"github.com/orrybaram/gary/internal/ui"
)

// Agent drives the conversation loop: read user input, call the model, run any
// tools it asks for, repeat.
type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
	toolset        []tools.Definition
}

// New builds an Agent from its dependencies. Nothing is read from the
// environment here; the caller owns configuration.
func New(
	client *anthropic.Client,
	getUserMessage func() (string, bool),
	toolset []tools.Definition,
) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		toolset:        toolset,
	}
}

func (a *Agent) Run(ctx context.Context) error {
	conversation := []anthropic.MessageParam{}

	fmt.Print(ui.Banner(Name))

	readUserInput := true

	for {
		if readUserInput {
			ui.UserPrompt()
			userInput, ok := a.getUserMessage()
			if !ok {
				break
			}

			// The API rejects empty text blocks, so ignore blank input.
			if strings.TrimSpace(userInput) == "" {
				continue
			}

			userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInput))
			conversation = append(conversation, userMessage)
		}

		spinner := ui.StartSpinner("Thinking")
		message, err := a.runInference(ctx, conversation)
		spinner.Stop()
		if err != nil {
			return err
		}
		// A truncated response can cut off mid-tool_use, which would otherwise
		// look like the model simply finished its turn. Say so loudly.
		if message.StopReason == anthropic.StopReasonMaxTokens {
			ui.Warn("response truncated at max_tokens (%d output tokens); ask me to continue",
				message.Usage.OutputTokens)
		}
		assistantMsg := stripEmptyText(message.ToParam())
		if len(assistantMsg.Content) > 0 {
			conversation = append(conversation, assistantMsg)
		}

		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, content := range message.Content {
			switch content.Type {
			case "text":
				if strings.TrimSpace(content.Text) == "" {
					continue
				}
				ui.AgentMessage(Name, content.Text)
			case "tool_use":
				result := a.executeTool(content.ID, content.Name, content.Input)
				toolResults = append(toolResults, result)
			}
		}
		if len(toolResults) == 0 {
			readUserInput = true
			continue
		}

		readUserInput = false
		conversation = append(conversation, anthropic.NewUserMessage(toolResults...))
	}

	return nil
}

// stripEmptyText removes empty/whitespace-only text blocks from an assistant
// message. The API returns 400 ("text content blocks must be non-empty") if
// such a block is echoed back in the conversation, which happens when the model
// emits an empty text block alongside tool_use.
func stripEmptyText(msg anthropic.MessageParam) anthropic.MessageParam {
	kept := make([]anthropic.ContentBlockParamUnion, 0, len(msg.Content))
	for _, block := range msg.Content {
		if block.OfText != nil && strings.TrimSpace(block.OfText.Text) == "" {
			continue
		}
		kept = append(kept, block)
	}
	msg.Content = kept
	return msg
}

func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	anthropicTools := []anthropic.ToolUnionParam{}
	for _, tool := range a.toolset {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeOpus5,
		MaxTokens: int64(8192),
		System:    []anthropic.TextBlockParam{{Text: systemPrompt}},
		Messages:  conversation,
		Tools:     anthropicTools,
	})

	return message, err
}

func (a *Agent) executeTool(id, name string, input json.RawMessage) anthropic.ContentBlockParamUnion {
	var toolDef tools.Definition
	var found bool
	for _, tool := range a.toolset {
		if tool.Name == name {
			toolDef = tool
			found = true
			break
		}
	}

	if !found {
		return anthropic.NewToolResultBlock(id, "tool not found", true)
	}

	ui.ToolCall(name, string(input))

	spinner := ui.StartSpinner(name)
	response, err := toolDef.Function(input)
	spinner.Stop()
	if err != nil {
		return anthropic.NewToolResultBlock(id, err.Error(), true)
	}
	return anthropic.NewToolResultBlock(id, response, false)
}
