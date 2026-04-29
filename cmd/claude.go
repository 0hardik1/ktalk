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
	APIEndpoint      = "https://api.anthropic.com/v1/messages"
	ModelName        = "claude-opus-4-7"
	AnthropicVersion = "2023-06-01"
	MaxTokens        = 1024
)

const systemPrompt = `You are a kubectl command generator. Convert the user's natural-language request into a single valid kubectl command.

Output rules:
- Return ONLY the command, with no explanation, prose, or markdown fences.
- The command MUST start with "kubectl".
- The command MUST be a single invocation: no pipes (|), no shell operators (&&, ||, ;), no command substitution. The runner rejects these.
- Use kubectl-native filtering and formatting: --sort-by, -o jsonpath, -o custom-columns, -o wide/yaml/json, --field-selector, --selector, --no-headers. Never rely on grep/awk/sed/wc/sort.
- Use --all-namespaces (or -A) when the request spans namespaces; otherwise respect the user's namespace if given.
- Quote and escape jsonpath/custom-columns values so the command runs as-is when copy-pasted.

Notes:
- Kubernetes has no built-in "users" resource. For service identities use ` + "`kubectl get serviceaccounts -A`" + `. For human access, inspect ` + "`clusterrolebindings`/`rolebindings`" + ` with -o custom-columns.
- For counts, prefer ` + "`kubectl get <resource> -A --no-headers`" + ` and let the caller count rows — do not pipe to wc.`

func formatAndValidateKubectlCommand(cmd string) (string, error) {
	cmd = strings.TrimSpace(cmd)
	cmd = strings.ReplaceAll(cmd, "```", "")
	cmd = strings.TrimSpace(cmd)

	if !strings.HasPrefix(cmd, "kubectl ") {
		return "", fmt.Errorf("invalid command: must start with 'kubectl'")
	}

	return cmd, nil
}

func ClaudeRequest(chatPrompt, apiKey string) error {
	client := resty.New()

	resp, err := client.R().
		SetHeader("x-api-key", apiKey).
		SetHeader("anthropic-version", AnthropicVersion).
		SetHeader("content-type", "application/json").
		SetBody(map[string]interface{}{
			"model":      ModelName,
			"max_tokens": MaxTokens,
			"system":     systemPrompt,
			"messages": []interface{}{
				map[string]interface{}{"role": "user", "content": chatPrompt},
			},
		}).
		Post(APIEndpoint)

	if err != nil {
		return fmt.Errorf("failed to send request to Anthropic API: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("Anthropic API returned non-200 status code: %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	body := resp.Body()
	if len(body) == 0 {
		return fmt.Errorf("received empty response from Anthropic API")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response from Anthropic API: %w, body: %s", err, string(body))
	}

	if respType, ok := response["type"].(string); ok && respType == "error" {
		if errObj, ok := response["error"].(map[string]interface{}); ok {
			if message, ok := errObj["message"].(string); ok {
				return fmt.Errorf("Anthropic API error: %s", message)
			}
		}
		return fmt.Errorf("Anthropic API error occurred")
	}

	contentField, ok := response["content"]
	if !ok {
		return fmt.Errorf("invalid response format: missing content field, response: %v", response)
	}

	contentArray, ok := contentField.([]interface{})
	if !ok {
		return fmt.Errorf("invalid response format: content is not an array, got: %T", contentField)
	}

	if len(contentArray) == 0 {
		return fmt.Errorf("invalid response format: empty content array")
	}

	block, ok := contentArray[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response format: content[0] is not an object, got: %T", contentArray[0])
	}

	text, ok := block["text"].(string)
	if !ok {
		return fmt.Errorf("invalid response format: content[0].text is not a string, got: %T", block["text"])
	}

	validatedCmd, err := formatAndValidateKubectlCommand(text)
	if err != nil {
		return fmt.Errorf("invalid command format: %w", err)
	}

	fmt.Println("Are you sure want to execute the following command? Press Enter to execute this: ", validatedCmd)

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
