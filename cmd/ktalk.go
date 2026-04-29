package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	genericclioptions.IOStreams
	AnthropicKey string
}

func NewPrompt(streams genericclioptions.IOStreams) *cobra.Command {
	o := &Options{
		IOStreams: streams,
	}

	// ktalk command definition
	cmd := &cobra.Command{
		Use:   "ktalk",
		Short: "ktalk talks to your Kubernetes cluster",
		Long:  "ktalk uses the Anthropic Claude API to generate kubectl commands based on natural language descriptions.\nNote: If your query ends with a question mark (?), you'll need to either quote your query, escape the question mark with a backslash, or use the special placeholder 'QUESTION' at the end.\n\nRunning 'kubectl ktalk' without arguments starts interactive mode.",
		Example: `  # Basic usage
  kubectl ktalk give me the list of containers in kube-system namespace
  
  # Using quotes for questions (recommended)
  kubectl ktalk "how many pods are running in the cluster?"
  
  # Escaping question marks
  kubectl ktalk how many pods are running in the cluster\?
  
  # Using the QUESTION placeholder
  kubectl ktalk how many users in the cluster QUESTION
  
  # Interactive mode
  kubectl ktalk`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.AnthropicKey = os.Getenv("ANTHROPIC_API_KEY")
			if o.AnthropicKey == "" {
				return fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
			}

			if len(args) == 0 {
				// No arguments provided, enter interactive mode
				return o.runInteractiveMode()
			}

			// Combine all arguments into a single prompt
			chatPrompt := strings.Join(args, " ")

			// Replace the QUESTION placeholder with an actual question mark
			if strings.HasSuffix(chatPrompt, " QUESTION") {
				chatPrompt = strings.TrimSuffix(chatPrompt, " QUESTION") + "?"
			}

			return o.run(chatPrompt)
		},
	}
	return cmd
}

func (o *Options) runInteractiveMode() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Entering interactive mode. Type 'exit' or 'quit' to exit.")

	for {
		fmt.Print("\nktalk> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		// Trim whitespace and newlines
		input = strings.TrimSpace(input)

		// Check for exit command
		if input == "exit" || input == "quit" {
			fmt.Println("Exiting interactive mode.")
			return nil
		}

		// Skip empty input
		if input == "" {
			continue
		}

		// Process the input
		if err := o.run(input); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			// Continue even if there's an error
		}
	}
}

func (o *Options) run(chatPrompt string) error {
	if err := ClaudeRequest(chatPrompt, o.AnthropicKey); err != nil {
		return fmt.Errorf("error when sending request to Anthropic API: %w", err)
	}
	return nil
}
