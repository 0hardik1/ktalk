package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewPrompt(streams genericclioptions.IOStreams) *cobra.Command {

	// ktalk command definition
	cmd := &cobra.Command{
		Use:          "ktalk",
		Short:        "ktalk talks to your Kubernetes cluster",
		Long:         "ktalk uses the OpenAI API to talk to generate kubectl commands for your Kubernetes cluster",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println("Please provide a message to ktalk.")
				fmt.Println("Usage: kubectl ktalk give me the list of containers in kube-system namespace")
				return nil
			}

			// Loop through all args and create a single string
			var ChatPrompt string
			for _, arg := range args {
				ChatPrompt += arg + " "
			}
			return run(ChatPrompt)

		},
	}
	return cmd
}

func run(ChatPrompt string) error {
	err := OpenAIRequest(ChatPrompt)
	if err != nil {
		fmt.Println("Error when sending request to OpenAI API", err)
	}

	return nil
}
