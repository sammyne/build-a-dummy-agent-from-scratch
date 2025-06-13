package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sammyne/build-a-dummy-agent-from-scratch/openai"
)

type Agent struct {
	httpClient     *http.Client // Use standard HTTP client
	getUserMessage func() (string, bool)
	tools          map[string]ToolDefinition // Use map for easy lookup by name
	model          string                    // Store the target model name
	systemPrompt   string                    // Store the system prompt
}

// --- Tool Definition (Mostly Unchanged, but Schema Generation Adapted) ---
type ToolDefinition struct {
	Name        string
	Description string
	// InputSchema now map[string]any to match OpenAI's parameter schema format
	InputSchema map[string]any
	Function    func(input json.RawMessage) (string, error) // Input is JSON string from OpenAI args
}

// Simplified flow within Agent.Run for basic chat
func (a *Agent) Run(ctx context.Context) error {
	// 初始化
	systemPromptMsg := openai.Message{
		Role:    "system",
		Content: a.systemPrompt,
	}
	conversation := []openai.Message{systemPromptMsg}

	for { // Outer loop for user input
		fmt.Print("\u001b[94mYou\u001b[0m: ") // Blue prompt for user
		userInput, ok := a.getUserMessage()
		if !ok { // Handle EOF or scanner error (e.g., ctrl-d)
			fmt.Println("\nExiting.")
			break
		}
		if userInput == "" {
			continue
		}

		// ... get userInput from console ...
		conversation = append(conversation, openai.Message{Role: "user", Content: userInput})

		// --- Call API ---
		resp, err := a.complete(ctx, conversation)
		if err != nil {
			fmt.Printf("\u001b[91mAPI Error\u001b[0m: %s\n", err.Error())
			continue // Let user try again
		}
		if len(resp.Choices) == 0 { /* handle no choices */
			continue
		}

		assistantMessage := resp.Choices[0].Message
		conversation = append(conversation, assistantMessage) // Add response to history

		// --- Print Text Response ---
		if assistantMessage.Content != "" {
			fmt.Printf("\u001b[93mAI\u001b[0m: %s\n", assistantMessage.Content)
		}

		// --- Tool Handling Logic would go here, but skipped for basic chat ---
		// In a basic chat without tool calls, the inner loop (if any) breaks immediately.

	} // End of outer loop

	return nil
}

func New(
	getUserMessage func() (string, bool),
	tools []ToolDefinition,
	model string,
) *Agent {
	toolMap := make(map[string]ToolDefinition)
	for _, tool := range tools {
		toolMap[tool.Name] = tool
	}

	return &Agent{
		httpClient:     &http.Client{Timeout: 60 * time.Second}, // Add a timeout
		getUserMessage: getUserMessage,
		tools:          toolMap,
		model:          model,
		// Define the system prompt here or pass it in
		systemPrompt: "You are a helpful Go programmer assistant. You have access to tools to interact with the local filesystem (read, list, edit files). Use them when appropriate to fulfill the user's request. When editing, be precise about the changes. Respond ONLY with tool calls if you need to use tools, otherwise respond with text.",
	}
}
