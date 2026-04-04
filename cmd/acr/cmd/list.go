package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/plexusone/agent-code-review/pkg/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List open pull requests",
	Long: `List open pull requests in a GitHub repository.

Example:
  acr list -o owner -r repo
  acr list -o owner -r repo --json`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("json", false, "Output as JSON")
}

func runList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	repoCfg, err := getRepoConfig(cmd)
	if err != nil {
		return err
	}

	client, err := config.CreateClient(ctx)
	if err != nil {
		return err
	}

	prs, err := client.ListOpenPRs(ctx, repoCfg.Owner, repoCfg.Repo)
	if err != nil {
		return fmt.Errorf("listing PRs: %w", err)
	}

	if len(prs) == 0 {
		fmt.Println("No open pull requests")
		return nil
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		data, err := json.MarshalIndent(prs, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	for _, pr := range prs {
		fmt.Printf("#%-6d %s (%s)\n", pr.Number, pr.Title, pr.Author)
	}

	return nil
}
