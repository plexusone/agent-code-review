package cmd

import (
	"context"
	"fmt"

	"github.com/plexusone/agent-code-review/pkg/config"
	"github.com/plexusone/agent-code-review/pkg/input"
	"github.com/plexusone/agent-code-review/pkg/review"
	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:   "comment <pr-number>",
	Short: "Add a general comment to a pull request",
	Long: `Add a general comment to a GitHub pull request.

This adds an issue comment (not a review) to the PR discussion.

Example:
  acr comment 123 -o owner -r repo -b "Thanks for the contribution!"
  acr comment 123 -o owner -r repo -f comment.md`,
	Args: cobra.ExactArgs(1),
	RunE: runComment,
}

func init() {
	rootCmd.AddCommand(commentCmd)
	commentCmd.Flags().StringP("body", "b", "", "Comment body text")
	commentCmd.Flags().StringP("file", "f", "", "Read comment body from file")
}

func runComment(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	repoCfg, err := getRepoConfig(cmd)
	if err != nil {
		return err
	}

	prNumber, err := input.ParsePRNumber(args[0])
	if err != nil {
		return err
	}

	body, err := getBody(cmd)
	if err != nil {
		return err
	}

	client, err := config.CreateClient(ctx)
	if err != nil {
		return err
	}

	result, err := client.CreateComment(ctx, &review.CommentInput{
		Owner:    repoCfg.Owner,
		Repo:     repoCfg.Repo,
		PRNumber: prNumber,
		Body:     body,
	})
	if err != nil {
		return fmt.Errorf("creating comment: %w", err)
	}

	fmt.Printf("Comment posted: %s\n", result.HTMLURL)
	return nil
}
