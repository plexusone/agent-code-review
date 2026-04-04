package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/plexusone/agent-code-review/pkg/config"
	"github.com/plexusone/agent-code-review/pkg/input"
	"github.com/plexusone/agent-code-review/pkg/review"
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review <pr-number>",
	Short: "Post a code review to a pull request",
	Long: `Post a code review to a GitHub pull request.

The review body can be provided via:
  - --body/-b flag
  - --file/-f flag (read from file)
  - stdin (if neither flag is provided)

Example:
  acr review 123 -o owner -r repo -e APPROVE -b "LGTM!"
  acr review 123 -o owner -r repo -e REQUEST_CHANGES -f review.md
  echo "Great work!" | acr review 123 -o owner -r repo -e COMMENT`,
	Args: cobra.ExactArgs(1),
	RunE: runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)
	reviewCmd.Flags().StringP("event", "e", "COMMENT", "Review event: APPROVE, REQUEST_CHANGES, or COMMENT")
	reviewCmd.Flags().StringP("body", "b", "", "Review body text")
	reviewCmd.Flags().StringP("file", "f", "", "Read review body from file")
}

func runReview(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	repoCfg, err := getRepoConfig(cmd)
	if err != nil {
		return err
	}

	prNumber, err := input.ParsePRNumber(args[0])
	if err != nil {
		return err
	}

	event, _ := cmd.Flags().GetString("event")
	event = strings.ToUpper(event)
	if err := input.ValidateReviewEvent(event); err != nil {
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

	result, err := client.CreateReview(ctx, &review.ReviewInput{
		Owner:    repoCfg.Owner,
		Repo:     repoCfg.Repo,
		PRNumber: prNumber,
		Event:    review.ReviewEvent(event),
		Body:     body,
	})
	if err != nil {
		return fmt.Errorf("creating review: %w", err)
	}

	fmt.Printf("Review posted: %s\n", result.HTMLURL)
	return nil
}

func getBody(cmd *cobra.Command) (string, error) {
	bodyText, _ := cmd.Flags().GetString("body")
	file, _ := cmd.Flags().GetString("file")

	// Check if stdin is a pipe
	stat, _ := os.Stdin.Stat()
	stdinIsPipe := (stat.Mode() & os.ModeCharDevice) == 0

	return input.ReadBody(input.BodySource{
		Text:        bodyText,
		File:        file,
		Stdin:       os.Stdin,
		StdinIsPipe: stdinIsPipe,
	})
}
