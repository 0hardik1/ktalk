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
)

func OpenAIRequest(ChatPrompt string) error {

	// Find API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := resty.New()

	kubectlCommandOnly := "give me just the command so that I can copy the command and paste it in my terminal, the command must start with kubectl without any quotes."

	// Create a new request to the OpenAI API
	resp, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "gpt-4-0125-preview",
			"messages": []interface{}{map[string]interface{}{"role": "system",
				"content": ChatPrompt + kubectlCommandOnly}},
			"max_tokens": 50,
		}).
		Post(APIEndpoint)

	if err != nil {
		fmt.Println("Error when sending request to OpenAI API", err)
		return err
	}

	body := resp.Body()

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error when unmarshalling response from OpenAI API", err)
		return err
	}

	// Extract the content from the JSON response
	content := response["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	fmt.Println("Are you sure want to execute the following command? Press Enter to execute this: ", content)

	// Agreed to execute?
	agreement, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
		return err
	}

	if agreement == "\n" {
		err := RunCommand(content)
		if err != nil {
			fmt.Println("Error when running command", err)
			return err
		}
		return nil
	}
	return nil
}
