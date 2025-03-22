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

// formatAndValidateKubectlCommand checks if the command is valid and properly formatted
func formatAndValidateKubectlCommand(cmd string) (string, error) {
	// Remove leading/trailing whitespace and newlines
	cmd = strings.TrimSpace(cmd)

	// Basic validation: must start with kubectl
	if !strings.HasPrefix(cmd, "kubectl ") {
		return "", fmt.Errorf("invalid command: must start with 'kubectl'")
	}

	// Remove any markdown code blocks if present
	cmd = strings.ReplaceAll(cmd, "```", "")
	cmd = strings.TrimSpace(cmd)

	return cmd, nil
}

func OpenAIRequest(chatPrompt, apiKey string) error {
	client := resty.New()

	kubectlCommandOnly := `
	You are a kubectl command generator. Your task is to convert natural language into valid kubectl commands.

	IMPORTANT: This system DOES NOT SUPPORT PIPES (|) in commands. Any command with a pipe will be rejected.

	Rules:
	1. Return ONLY the kubectl command without any explanations or markdown
	2. The command must start with 'kubectl'
	3. NEVER use pipes (|) - they are not supported and will cause your command to fail
	4. Use appropriate output formats based on the query (-o wide, -o yaml, -o json, -o custom-columns)
	5. Use --all-namespaces when the query involves looking across all namespaces
	6. Keep commands simple, readable, and efficient
	7. For security context related queries, consider both pod-level and container-level settings
	8. Use jsonpath or custom-columns for extracting specific fields when needed
	9. Ensure all quotes, brackets, and braces are properly closed and escaped
	10. Generate commands that can be copy-pasted and executed without any modifications
	11. For sorting resources, ONLY use kubectl's --sort-by option (not external sort commands)
	12. Do not use grep, awk, sed or any other filtering that requires pipes
	13. For listing containers, use: kubectl get pods --all-namespaces -o=custom-columns="NAMESPACE:.metadata.namespace,POD:.metadata.name,CONTAINER:.spec.containers[*].name"
	14. For counting resources, use either:
	    - kubectl get [resource] --no-headers | wc -l (THIS WON'T WORK - has a pipe)
	    - kubectl get [resource] -o name | wc -l (THIS WON'T WORK - has a pipe)
	    Instead use: kubectl get [resource] --no-headers
	15. For users specifically, try: kubectl get serviceaccounts --all-namespaces (Kubernetes doesn't have a built-in "users" resource type)
	16. To find human users with access, use: kubectl get clusterrolebindings -o=custom-columns="NAME:.metadata.name,ROLE:.roleRef.name,SUBJECTS:.subjects[*].name"
	`

	// Create a new request to the OpenAI API
	resp, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": ModelName,
			"messages": []interface{}{map[string]interface{}{"role": "system",
				"content": kubectlCommandOnly}, map[string]interface{}{"role": "user",
				"content": chatPrompt}},
			"max_tokens": 150,
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

	// Validate and format the command
	validatedCmd, err := formatAndValidateKubectlCommand(contentStr)
	if err != nil {
		return fmt.Errorf("invalid command format: %w", err)
	}

	fmt.Println("Are you sure want to execute the following command? Press Enter to execute this: ", validatedCmd)

	// Agreed to execute?
	agreement, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	if agreement == "\n" {
		if err := RunCommand(validatedCmd); err != nil {
			return fmt.Errorf("failed to run command: %w", err)
		}
	}
	return nil
}
