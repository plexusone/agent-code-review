package cmd

import (
	"context"
	"fmt"

	"github.com/plexusone/agent-code-review/pkg/config"
	"github.com/plexusone/agent-code-review/pkg/input"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <pr-number>",
	Short: "Get the diff for a pull request",
	Long: `Fetch and display the diff for a GitHub pull request.

Example:
  acr diff 123 -o owner -r repo`,
	Args: cobra.ExactArgs(1),
	RunE: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	repoCfg, err := getRepoConfig(cmd)
	if err != nil {
		return err
	}

	prNumber, err := input.ParsePRNumber(args[0])
	if err != nil {
		return err
	}

	client, err := config.CreateClient(ctx)
	if err != nil {
		return err
	}

	diff, err := client.GetPRDiff(ctx, repoCfg.Owner, repoCfg.Repo, prNumber)
	if err != nil {
		return fmt.Errorf("getting diff: %w", err)
	}

	fmt.Print(diff)
	return nil
}
