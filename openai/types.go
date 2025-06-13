package openai

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`       // The assistant's response message
	FinishReason string  `json:"finish_reason"` // e.g., "stop", "tool_calls"
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // Arguments are a *string* containing JSON
}

type FunctionDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters"` // Use map[string]any to represent JSON schema object
}

type Message struct {
	Role       string     `json:"role"`                   // "system", "user", "assistant", "tool"
	Content    string     `json:"content,omitempty"`      // For text content or tool result
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // For assistant requesting tools
	ToolCallID string     `json:"tool_call_id,omitempty"` // For tool role messages
	Name       string     `json:"name,omitempty"`         // For tool role messages (function name) - Optional by OpenAI spec but sometimes useful
}

type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Tools       []Tool    `json:"tools,omitempty"`
	ToolChoice  any       `json:"tool_choice,omitempty"` // "auto" or specific tool
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
	// Add other OpenAI parameters as needed (top_p, stream, etc.)
}

type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Tool struct {
	Type     string             `json:"type"` // Always "function" for now
	Function FunctionDefinition `json:"function"`
}

type ToolCall struct {
	ID       string       `json:"id"`   // ID to match with tool response
	Type     string       `json:"type"` // Always "function"
	Function FunctionCall `json:"function"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
