package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	// APIEndpoint is the URL for the OpenAI API
	APIEndpoint = "https://api.openai.com/v1/chat/completions"
	// Using gpt-3.5-turbo as it's more stable, cost-effective, and sufficient for kubectl command generation
	ModelName = "gpt-3.5-turbo"
)

func OpenAIRequest(chatPrompt, apiKey string) error {
	client := resty.New()

	kubectlCommandOnly := `
	Response Requirements: 
	1. The response should be a just a kubectl command. 
	2. The command must start with kubectl without any quotes.`

	// Create a new request to the OpenAI API
	resp, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": ModelName,
			"messages": []interface{}{map[string]interface{}{"role": "system",
				"content": chatPrompt + kubectlCommandOnly}},
			"max_tokens": 50,
		}).
		Post(APIEndpoint)

	if err != nil {
		return fmt.Errorf("failed to send request to OpenAI API: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode() != 200 {
		return fmt.Errorf("OpenAI API returned non-200 status code: %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	body := resp.Body()
	if len(body) == 0 {
		return fmt.Errorf("received empty response from OpenAI API")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response from OpenAI API: %w, body: %s", err, string(body))
	}

	// Check for error in OpenAI response
	if err, ok := response["error"].(map[string]interface{}); ok {
		if message, ok := err["message"].(string); ok {
			return fmt.Errorf("OpenAI API error: %s", message)
		}
		return fmt.Errorf("OpenAI API error occurred")
	}

	// Safely extract the content from the JSON response
	choices, ok := response["choices"]
	if !ok {
		return fmt.Errorf("invalid response format: missing choices field, response: %v", response)
	}

	choicesArray, ok := choices.([]interface{})
	if !ok {
		return fmt.Errorf("invalid response format: choices is not an array, got: %T", choices)
	}

	if len(choicesArray) == 0 {
		return fmt.Errorf("invalid response format: empty choices array")
	}

	choice, ok := choicesArray[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response format: choice is not an object, got: %T", choicesArray[0])
	}

	message, ok := choice["message"]
	if !ok {
		return fmt.Errorf("invalid response format: missing message field in choice")
	}

	messageMap, ok := message.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response format: message is not an object, got: %T", message)
	}

	content, ok := messageMap["content"]
	if !ok {
		return fmt.Errorf("invalid response format: missing content field in message")
	}

	contentStr, ok := content.(string)
	if !ok {
		return fmt.Errorf("invalid response format: content is not a string, got: %T", content)
	}

	// Clean up the command
	contentStr = strings.TrimSpace(contentStr)
	if !strings.HasPrefix(contentStr, "kubectl ") {
		return fmt.Errorf("invalid command format: command must start with 'kubectl', got: %s", contentStr)
	}

	fmt.Println("Are you sure want to execute the following command? Press Enter to execute this: ", contentStr)

	// Agreed to execute?
	agreement, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	if agreement == "\n" {
		if err := RunCommand(contentStr); err != nil {
			return fmt.Errorf("failed to run command: %w", err)
		}
	}
	return nil
}
