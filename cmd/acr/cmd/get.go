package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <pr-number>",
	Short: "Get pull request details",
	Long: `Fetch and display details for a GitHub pull request.

Example:
  acr get 123 -o owner -r repo
  acr get 123 -o owner -r repo --json`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().Bool("json", false, "Output as JSON")
}

func runGet(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	owner, repo, err := getOwnerRepo(cmd)
	if err != nil {
		return err
	}

	prNumber, err := parsePRNumber(args[0])
	if err != nil {
		return err
	}

	client, err := createClient(ctx)
	if err != nil {
		return err
	}

	pr, err := client.GetPR(ctx, owner, repo, prNumber)
	if err != nil {
		return fmt.Errorf("getting PR: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		data, err := json.MarshalIndent(pr, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("PR #%d: %s\n", pr.Number, pr.Title)
	fmt.Printf("Author: %s\n", pr.Author)
	fmt.Printf("State: %s\n", pr.State)
	fmt.Printf("Branch: %s -> %s\n", pr.Head, pr.Base)
	fmt.Printf("Commits: %d\n", pr.Commits)
	fmt.Printf("URL: %s\n", pr.HTMLURL)
	if pr.Body != "" {
		fmt.Printf("\n%s\n", pr.Body)
	}

	return nil
}
