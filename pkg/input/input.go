// Package input provides input parsing utilities for agent-code-review.
package input

import (
	"fmt"
	"io"
	"os"
)

// ParsePRNumber parses a PR number from a string.
func ParsePRNumber(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("invalid PR number: %s", s)
	}
	// Check that all characters are digits
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid PR number: %s", s)
		}
	}
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid PR number: %s", s)
	}
	return n, nil
}

// BodySource specifies where to read the body from.
type BodySource struct {
	// Text is the body text provided directly.
	Text string
	// File is the path to read the body from.
	File string
	// Stdin is the reader for stdin input (nil to skip stdin).
	Stdin io.Reader
	// StdinIsPipe indicates whether stdin is a pipe (not a terminal).
	StdinIsPipe bool
}

// ReadBody reads the body from the specified source.
// Priority: Text > File > Stdin
func ReadBody(src BodySource) (string, error) {
	if src.Text != "" {
		return src.Text, nil
	}

	if src.File != "" {
		data, err := os.ReadFile(src.File)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return string(data), nil
	}

	if src.Stdin != nil && src.StdinIsPipe {
		data, err := io.ReadAll(src.Stdin)
		if err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		return string(data), nil
	}

	return "", fmt.Errorf("body required: provide text, file path, or pipe to stdin")
}

// ValidateReviewEvent validates a review event string.
func ValidateReviewEvent(event string) error {
	switch event {
	case "APPROVE", "REQUEST_CHANGES", "COMMENT":
		return nil
	default:
		return fmt.Errorf("invalid event: %s (must be APPROVE, REQUEST_CHANGES, or COMMENT)", event)
	}
}
