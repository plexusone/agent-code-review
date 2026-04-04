// Package main provides the entry point for the agent-code-review CLI.
package main

import (
	"os"

	"github.com/plexusone/agent-code-review/cmd/acr/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
