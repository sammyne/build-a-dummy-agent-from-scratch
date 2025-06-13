package openai

import "os"

// Read from environment variables
var (
	APIKey      = os.Getenv("OPENAI_API_KEY")                           // Use OPENAI_API_KEY now
	APIEndpoint = os.Getenv("OPENAI_API_BASE") + "/v1/chat/completions" // Allow overriding base URL
	Model       = os.Getenv("OPENAI_MODEL")                             // Allow specifying model
)
