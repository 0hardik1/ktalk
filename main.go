package main

import (
	"os"

	"github.com/0hardik1/ktalk/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	prompt := cmd.NewPrompt(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := prompt.Execute(); err != nil {
		os.Exit(1)
	}
}
