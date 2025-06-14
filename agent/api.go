package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sammyne/build-a-dummy-agent-from-scratch/openai"
	"github.com/sammyne/build-a-dummy-agent-from-scratch/tools"
)

type Agent struct {
	httpClient     *http.Client // Use standard HTTP client
	getUserMessage func() (string, bool)
	tools          map[string]tools.Definition // Use map for easy lookup by name
	model          string                      // Store the target model name
	systemPrompt   string                      // Store the system prompt
}

// Simplified flow within Agent.Run for basic chat
func (a *Agent) Run(ctx context.Context) error {
	// 初始化
	systemPromptMsg := openai.Message{
		Role:    "system",
		Content: a.systemPrompt,
	}
	conversation := []openai.Message{systemPromptMsg}

	fmt.Println("Chat with AI (use 'ctrl-c' to quit)")

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

		// --- Main Loop: Call API -> Handle Response -> Execute Tools -> Call API ---
		for { // Inner loop to handle potential multi-turn tool calls
			// Call OpenAI API
			conversationJSON,_:=json.MarshalIndent(conversation,"","  ")
			fmt.Printf("conversations: %s\n",conversationJSON)

			resp, err := a.complete(ctx, conversation)
			if err != nil {
				fmt.Printf("\u001b[91mAPI Error\u001b[0m: %s\n", err.Error())
				continue // Let user try again
			}
			if len(resp.Choices) == 0 { /* handle no choices */
				fmt.Println("\u001b[91mError\u001b[0m: OpenAI response contained no choices.")
				continue
			}

			assistantMessage := resp.Choices[0].Message
			// Add assistant's message (text and/or tool calls) to conversation
			conversation = append(conversation, assistantMessage) // Add response to history

			// --- Print Text Response ---
			if assistantMessage.Content != "" {
				fmt.Printf("\u001b[93mAI\u001b[0m: %s\n", assistantMessage.Content)
			}

			// --- Handle Tool Calls ---
			if len(assistantMessage.ToolCalls) == 0 {
				// No tools called, break inner loop and wait for next user input
				break
			}

			// In a basic chat without tool calls, the inner loop (if any) breaks immediately.
			toolOutputs := a.callTools(assistantMessage.ToolCalls)
			// Add all tool results to conversation history
			conversation = append(conversation, toolOutputs...)
		}

	} // End of outer loop

	return nil
}

func (a *Agent) callTools(calls []openai.ToolCall) []openai.Message {
	if len(calls) == 0 {
		return nil
	}

	var out []openai.Message
	for _, c := range calls {
		if c.Type != "function" {
			continue // Skip non-function tool calls if any
		}

		toolName := c.Function.Name
		toolArgs := c.Function.Arguments // This is a JSON *string*

		fmt.Printf("\u001b[92mTool Call\u001b[0m: %s(%s)\n", toolName, toolArgs) // Green

		resultMsg := openai.Message{
			Role:       "tool",
			ToolCallID: c.ID,
			Name:       toolName,
		}

		toolDef, found := a.tools[toolName]
		if !found {
			errorMsg := fmt.Sprintf("tool '%s' not found by agent", toolName)
			fmt.Printf("\u001b[91mTool Error\u001b[0m: %s\n", errorMsg)
			resultMsg.Content = errorMsg // Report error back to OpenAI
		} else {
			// Execute the actual tool function
			// Note: toolArgs is a JSON string, pass it as json.RawMessage
			toolOutput, err := toolDef.Function(json.RawMessage(toolArgs))
			if err != nil {
				errorMsg := fmt.Sprintf("error executing tool '%s': %s", toolName, err.Error())
				fmt.Printf("\u001b[91mTool Error\u001b[0m: %s\n", errorMsg)
				resultMsg.Content = errorMsg // Report error back to OpenAI
			} else {
				// Log successful tool execution result (optional)
				// fmt.Printf("\u001b[92mTool Result\u001b[0m: %s\n", toolOutput)
				resultMsg.Content = toolOutput // Send success result back to OpenAI
			}
		}
		out = append(out, resultMsg)
	} // End of processing tool calls for one response

	return out
}

func New(
	getUserMessage func() (string, bool),
	toolList []tools.Definition,
	model string,
) *Agent {
	toolMap := make(map[string]tools.Definition)
	for _, tool := range toolList {
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
