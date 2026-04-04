// Package cmd provides the CLI commands for agent-code-review.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "acr",
	Short: "Agent Code Review - AI-powered code review for GitHub PRs",
	Long: `Agent Code Review (acr) is an AI-powered code review tool for GitHub Pull Requests.

It can be used as a CLI tool or as an MCP server for integration with Claude Code.

Reviews are posted as a GitHub App, appearing as "PlexusOne Code Review[bot]",
clearly distinguishing AI-assisted reviews from human reviews.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("owner", "o", "", "Repository owner (user or org)")
	rootCmd.PersistentFlags().StringP("repo", "r", "", "Repository name")
}

func getOwnerRepo(cmd *cobra.Command) (string, string, error) {
	owner, _ := cmd.Flags().GetString("owner")
	repo, _ := cmd.Flags().GetString("repo")

	if owner == "" {
		owner = os.Getenv("GITHUB_OWNER")
	}
	if repo == "" {
		repo = os.Getenv("GITHUB_REPO")
	}

	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("owner and repo are required (use --owner/-o and --repo/-r flags or GITHUB_OWNER/GITHUB_REPO env vars)")
	}

	return owner, repo, nil
}
