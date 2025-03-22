package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

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

	body := resp.Body()

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response from OpenAI API: %w", err)
	}

	// Extract the content from the JSON response
	content, ok := response["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return fmt.Errorf("unexpected response format from OpenAI API")
	}

	fmt.Println("Are you sure want to execute the following command? Press Enter to execute this: ", content)

	// Agreed to execute?
	agreement, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	if agreement == "\n" {
		if err := RunCommand(content); err != nil {
			return fmt.Errorf("failed to run command: %w", err)
		}
	}
	return nil
}
