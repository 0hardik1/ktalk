package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	genericclioptions.IOStreams
	OpenAIKey string
}

func NewPrompt(streams genericclioptions.IOStreams) *cobra.Command {
	o := &Options{
		IOStreams: streams,
	}

	// ktalk command definition
	cmd := &cobra.Command{
		Use:          "ktalk",
		Short:        "ktalk talks to your Kubernetes cluster",
		Long:         "ktalk uses the OpenAI API to generate kubectl commands based on natural language descriptions",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("please provide a message to ktalk\nUsage: kubectl ktalk give me the list of containers in kube-system namespace")
			}

			// Get OpenAI API key from environment
			o.OpenAIKey = os.Getenv("OPENAI_API_KEY")
			if o.OpenAIKey == "" {
				return fmt.Errorf("OPENAI_API_KEY environment variable is not set")
			}

			// Combine all arguments into a single prompt
			chatPrompt := strings.Join(args, " ")
			return o.run(chatPrompt)
		},
	}
	return cmd
}

func (o *Options) run(chatPrompt string) error {
	if err := OpenAIRequest(chatPrompt, o.OpenAIKey); err != nil {
		return fmt.Errorf("error when sending request to OpenAI API: %w", err)
	}
	return nil
}
