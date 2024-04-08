package cmd

import (
	"os"
	"os/exec"
)

// function that takes a shell command as input and returns the output of the command
func RunCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
