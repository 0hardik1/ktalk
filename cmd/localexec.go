package cmd

import (
	"context"
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

// RunCommand executes a kubectl command with proper security and output handling
func RunCommand(command string) error {
	// Validate the command first
	if err := validateKubectlCommand(command); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Split command into args
	args := strings.Fields(command)
	if len(args) < 2 {
		return fmt.Errorf("invalid command format")
	}

	// Create command with context
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	// Set up command IO
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Create a channel to handle command interruption
	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Wait for command completion or interruption
	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Wait()
	}()

	select {
	case err := <-errChan:
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return fmt.Errorf("command failed with exit code %d: %s", exitErr.ExitCode(), exitErr.Stderr)
			}
			return fmt.Errorf("command failed: %w", err)
		}
	case <-ctx.Done():
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill command: %w", err)
		}
		return fmt.Errorf("command timed out after %v", defaultTimeout)
	}

	return nil
}
