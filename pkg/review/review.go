// Package review provides a high-level API for GitHub code reviews.
package review

import (
	"context"
	"fmt"

	"github.com/google/go-github/v84/github"
	"github.com/grokify/gogithub/auth"
	"github.com/grokify/gogithub/pr"
)

// ReviewFooter is appended to all reviews for transparency.
const ReviewFooter = "\n\n---\n<sub>🤖 Powered by Claude • PlexusOne Code Review</sub>"

// Client provides code review operations.
type Client struct {
	gh *github.Client
}

// NewClient creates a new review client from a GitHub client.
func NewClient(gh *github.Client) *Client {
	return &Client{gh: gh}
}

// NewClientFromAppConfig creates a new review client using GitHub App authentication.
func NewClientFromAppConfig(ctx context.Context, cfg *auth.AppConfig) (*Client, error) {
	gh, err := auth.NewAppClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating github client: %w", err)
	}
	return &Client{gh: gh}, nil
}

// NewClientFromToken creates a new review client using a personal access token.
func NewClientFromToken(ctx context.Context, token string) *Client {
	gh := auth.NewGitHubClient(ctx, token)
	return &Client{gh: gh}
}

// ReviewEvent represents the type of review action.
type ReviewEvent string

const (
	// EventApprove approves the pull request.
	EventApprove ReviewEvent = "APPROVE"
	// EventRequestChanges requests changes to the pull request.
	EventRequestChanges ReviewEvent = "REQUEST_CHANGES"
	// EventComment adds a review comment without approval or rejection.
	EventComment ReviewEvent = "COMMENT"
)

// ReviewInput contains parameters for creating a review.
type ReviewInput struct {
	Owner    string
	Repo     string
	PRNumber int
	Event    ReviewEvent
	Body     string
}

// ReviewResult contains the result of a review operation.
type ReviewResult struct {
	ID      int64
	HTMLURL string
}

// CreateReview posts a code review to a pull request.
// The review footer is automatically appended to the body.
func (c *Client) CreateReview(ctx context.Context, input *ReviewInput) (*ReviewResult, error) {
	body := input.Body + ReviewFooter
	review, err := pr.CreateReview(ctx, c.gh, input.Owner, input.Repo, input.PRNumber, pr.ReviewEvent(input.Event), body)
	if err != nil {
		return nil, fmt.Errorf("creating review: %w", err)
	}
	return &ReviewResult{
		ID:      review.GetID(),
		HTMLURL: review.GetHTMLURL(),
	}, nil
}

// CommentInput contains parameters for creating a PR comment.
type CommentInput struct {
	Owner    string
	Repo     string
	PRNumber int
	Body     string
}

// CommentResult contains the result of a comment operation.
type CommentResult struct {
	ID      int64
	HTMLURL string
}

// CreateComment adds a general comment to a pull request.
// The review footer is automatically appended to the body.
func (c *Client) CreateComment(ctx context.Context, input *CommentInput) (*CommentResult, error) {
	body := input.Body + ReviewFooter
	comment, err := pr.CreateIssueComment(ctx, c.gh, input.Owner, input.Repo, input.PRNumber, body)
	if err != nil {
		return nil, fmt.Errorf("creating comment: %w", err)
	}
	return &CommentResult{
		ID:      comment.GetID(),
		HTMLURL: comment.GetHTMLURL(),
	}, nil
}

// LineCommentInput contains parameters for creating a line comment.
type LineCommentInput struct {
	Owner    string
	Repo     string
	PRNumber int
	CommitID string
	Path     string
	Line     int
	Body     string
}

// CreateLineComment adds a comment on a specific line in a PR diff.
func (c *Client) CreateLineComment(ctx context.Context, input *LineCommentInput) (*CommentResult, error) {
	comment, err := pr.CreateLineComment(ctx, c.gh, input.Owner, input.Repo, input.PRNumber, input.CommitID, input.Path, input.Body, input.Line)
	if err != nil {
		return nil, fmt.Errorf("creating line comment: %w", err)
	}
	return &CommentResult{
		ID:      comment.GetID(),
		HTMLURL: comment.GetHTMLURL(),
	}, nil
}

// PRInfo contains pull request information.
type PRInfo struct {
	Number  int
	Title   string
	Body    string
	State   string
	Author  string
	Head    string
	Base    string
	Commits int
	HTMLURL string
}

// GetPR retrieves pull request details.
func (c *Client) GetPR(ctx context.Context, owner, repo string, number int) (*PRInfo, error) {
	ghPR, err := pr.GetPR(ctx, c.gh, owner, repo, number)
	if err != nil {
		return nil, err
	}
	return &PRInfo{
		Number:  ghPR.GetNumber(),
		Title:   ghPR.GetTitle(),
		Body:    ghPR.GetBody(),
		State:   ghPR.GetState(),
		Author:  ghPR.GetUser().GetLogin(),
		Head:    ghPR.GetHead().GetRef(),
		Base:    ghPR.GetBase().GetRef(),
		Commits: ghPR.GetCommits(),
		HTMLURL: ghPR.GetHTMLURL(),
	}, nil
}

// GetPRDiff retrieves the diff for a pull request.
func (c *Client) GetPRDiff(ctx context.Context, owner, repo string, number int) (string, error) {
	return pr.GetPRDiff(ctx, c.gh, owner, repo, number)
}

// PRSummary contains basic pull request information for listings.
type PRSummary struct {
	Number  int
	Title   string
	Author  string
	Head    string
	HTMLURL string
}

// ListOpenPRs lists open pull requests in a repository.
func (c *Client) ListOpenPRs(ctx context.Context, owner, repo string) ([]PRSummary, error) {
	prs, err := pr.ListPRs(ctx, c.gh, owner, repo, &github.PullRequestListOptions{
		State: "open",
		ListOptions: github.ListOptions{
			PerPage: 30,
		},
	})
	if err != nil {
		return nil, err
	}

	result := make([]PRSummary, len(prs))
	for i, p := range prs {
		result[i] = PRSummary{
			Number:  p.GetNumber(),
			Title:   p.GetTitle(),
			Author:  p.GetUser().GetLogin(),
			Head:    p.GetHead().GetRef(),
			HTMLURL: p.GetHTMLURL(),
		}
	}
	return result, nil
}

// Approve approves a pull request with an optional comment.
func (c *Client) Approve(ctx context.Context, owner, repo string, number int, body string) (*ReviewResult, error) {
	return c.CreateReview(ctx, &ReviewInput{
		Owner:    owner,
		Repo:     repo,
		PRNumber: number,
		Event:    EventApprove,
		Body:     body,
	})
}

// RequestChanges requests changes on a pull request.
func (c *Client) RequestChanges(ctx context.Context, owner, repo string, number int, body string) (*ReviewResult, error) {
	return c.CreateReview(ctx, &ReviewInput{
		Owner:    owner,
		Repo:     repo,
		PRNumber: number,
		Event:    EventRequestChanges,
		Body:     body,
	})
}

// Comment adds a review comment without approval or rejection.
func (c *Client) Comment(ctx context.Context, owner, repo string, number int, body string) (*ReviewResult, error) {
	return c.CreateReview(ctx, &ReviewInput{
		Owner:    owner,
		Repo:     repo,
		PRNumber: number,
		Event:    EventComment,
		Body:     body,
	})
}
