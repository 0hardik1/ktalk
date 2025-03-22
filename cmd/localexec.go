package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// Default timeout for command execution
	defaultTimeout = 5 * time.Minute
)

// validateKubectlCommand checks if the command is a valid kubectl command
func validateKubectlCommand(command string) error {
	// Remove leading/trailing whitespace
	command = strings.TrimSpace(command)

	// Check if command starts with kubectl
	if !strings.HasPrefix(command, "kubectl ") {
		return fmt.Errorf("invalid command: must start with 'kubectl'")
	}

	// Basic security checks
	if strings.Contains(command, "&&") || strings.Contains(command, "||") || strings.Contains(command, ";") {
		return fmt.Errorf("invalid command: contains shell operators")
	}

	return nil
}

// RunCommand executes the given kubectl command
func RunCommand(command string) error {
	// Validate the command - no pipes with kubectl
	if strings.Contains(command, "|") {
		return fmt.Errorf("kubectl commands with pipes are not supported in this version. Try using an alternative command without pipes")
	}

	// Use shell to execute the command
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed with exit code %d: %s", cmd.ProcessState.ExitCode(), err)
	}
	return nil
}
