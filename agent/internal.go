package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sammyne/build-a-dummy-agent-from-scratch/openai"
)

// complete uses standard library http client
func (a *Agent) complete(ctx context.Context, conversation []openai.Message) (*openai.Response, error) {
	// Prepare tools in OpenAI format
	tools := []openai.Tool{}

	for _, toolDef := range a.tools {
		tool := openai.Tool{
			Type: "function",
			Function: openai.FunctionDefinition{
				Name:        toolDef.Name,
				Description: toolDef.Description,
				Parameters:  toolDef.InputSchema,
			},
		}

		tools = append(tools, tool)
	}

	// Build request payload
	request := openai.Request{
		Model:       a.model,
		Messages:    conversation,
		Tools:       tools,
		ToolChoice:  "auto", // Let the model decide when to use tools
		MaxTokens:   2048,   // Or make configurable
		Temperature: 0.7,    // Reasonable default
	}

	// Marshal payload to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}
	fmt.Printf("request: %s\n", requestJSON)
	fmt.Printf("api-key: %s\n", openai.APIKey)
	fmt.Printf("endpoint: %s\n", openai.APIEndpoint)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", openai.APIEndpoint, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openai.APIKey)
	// Add other headers if required by the specific OpenAI-compatible provider

	// Send request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		// Try to include error message from response body
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Unmarshal response JSON
	var response openai.Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		// Include raw body in error for debugging
		return nil, fmt.Errorf("failed to unmarshal response JSON: %w. Body: %s", err, string(bodyBytes))
	}

	return &response, nil
}
