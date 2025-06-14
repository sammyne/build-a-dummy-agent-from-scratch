package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/sammyne/build-a-dummy-agent-from-scratch/agent"
	"github.com/sammyne/build-a-dummy-agent-from-scratch/openai"
	"github.com/sammyne/build-a-dummy-agent-from-scratch/tools"
)

func main() {
	validateEnv()

	// --- Setup Input ---
	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "\u001b[91mError reading input: %v\u001b[0m\n", err)
				return "", false
			}
			return "", false // EOF
		}
		return scanner.Text(), true
	}

	tools := []tools.Definition{
		tools.EditFileDefinition,
		tools.ListFilesDefinition,
		tools.ReadFileDefinition,
	}

	// --- Create and Run Agent ---
	agent := agent.New(getUserMessage, tools, openai.Model)
	if err := agent.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "\u001b[91mAgent exited with error: %s\u001b[0m\n", err.Error())
		os.Exit(1)
	}
}

const DEFAULT_API_ENDPOINT = "https://api.siliconflow.cn/v1/chat/completions"

func validateEnv() {
	// --- Configuration Checks ---
	if openai.APIKey == "" {
		fmt.Fprintln(os.Stderr, "\u001b[91mError: OPENAI_API_KEY environment variable not set.\u001b[0m")
		os.Exit(1)
	}

	if os.Getenv("OPENAI_API_ENDPOINT") == "" {
		// Default to official OpenAI endpoint if base URL not set
		openai.APIEndpoint = DEFAULT_API_ENDPOINT
		fmt.Printf("Info: OPENAI_API_ENDPOINT not set, defaulting to %s\n", DEFAULT_API_ENDPOINT)
	}
	if openai.Model == "" {
		// Default model if not set
		openai.Model = "deepseek-ai/DeepSeek-V3" // 或其他模型
		fmt.Printf("Info: OPENAI_MODEL not set, defaulting to %s\n", openai.Model)
	}
}
