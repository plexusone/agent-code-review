package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/grokify/gogithub/auth"
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

	owner, repo, err := getOwnerRepo(cmd)
	if err != nil {
		return err
	}

	prNumber, err := parsePRNumber(args[0])
	if err != nil {
		return err
	}

	event, _ := cmd.Flags().GetString("event")
	event = strings.ToUpper(event)
	if event != "APPROVE" && event != "REQUEST_CHANGES" && event != "COMMENT" {
		return fmt.Errorf("invalid event: %s (must be APPROVE, REQUEST_CHANGES, or COMMENT)", event)
	}

	body, err := getBody(cmd)
	if err != nil {
		return err
	}

	client, err := createClient(ctx)
	if err != nil {
		return err
	}

	result, err := client.CreateReview(ctx, &review.ReviewInput{
		Owner:    owner,
		Repo:     repo,
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

func parsePRNumber(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid PR number: %s", s)
	}
	return n, nil
}

func getBody(cmd *cobra.Command) (string, error) {
	body, _ := cmd.Flags().GetString("body")
	if body != "" {
		return body, nil
	}

	file, _ := cmd.Flags().GetString("file")
	if file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		return string(data), nil
	}

	// Read from stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		return string(data), nil
	}

	return "", fmt.Errorf("review body required: use --body, --file, or pipe to stdin")
}

func createClient(ctx context.Context) (*review.Client, error) {
	// Try GitHub App auth first
	cfg, err := auth.LoadAppConfig()
	if err == nil {
		return review.NewClientFromAppConfig(ctx, cfg)
	}

	// Fall back to token auth
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("no authentication configured: set up GitHub App config or GITHUB_TOKEN env var")
	}

	return review.NewClientFromToken(ctx, token), nil
}
