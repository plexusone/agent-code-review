// Package review provides a high-level API for GitHub code reviews.
package review

import (
	"context"
	"fmt"

	"github.com/grokify/gogithub"
	"github.com/grokify/gogithub/auth"
	"github.com/grokify/gogithub/clientv1"
)

// ReviewFooter is appended to all reviews for transparency.
const ReviewFooter = "\n\n---\n<sub>🤖 Powered by Claude • PlexusOne Code Review</sub>"

// Client provides code review operations.
type Client struct {
	client clientv1.Client
}

// NewClient creates a new review client from a version-isolated GitHub client.
func NewClient(client clientv1.Client) *Client {
	return &Client{client: client}
}

// NewClientFromAppConfig creates a new review client using GitHub App authentication.
func NewClientFromAppConfig(ctx context.Context, cfg *auth.AppConfig) (*Client, error) {
	gh, err := auth.NewAppClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating github client: %w", err)
	}
	return &Client{client: clientv1.NewClientFromRaw(gh)}, nil
}

// NewClientFromToken creates a new review client using a personal access token.
func NewClientFromToken(ctx context.Context, token string) (*Client, error) {
	client, err := clientv1.NewClient(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("creating github client: %w", err)
	}
	return &Client{client: client}, nil
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
	review, err := c.client.CreatePullRequestReview(ctx, input.Owner, input.Repo, input.PRNumber, &clientv1.CreateReviewInput{
		Event: string(input.Event),
		Body:  body,
	})
	if err != nil {
		return nil, fmt.Errorf("creating review: %w", err)
	}
	return &ReviewResult{
		ID:      review.ID,
		HTMLURL: review.HTMLURL,
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
	comment, err := c.client.CreateIssueComment(ctx, input.Owner, input.Repo, input.PRNumber, body)
	if err != nil {
		return nil, fmt.Errorf("creating comment: %w", err)
	}
	return &CommentResult{
		ID:      comment.ID,
		HTMLURL: comment.HTMLURL,
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
	comment, err := c.client.CreatePullRequestComment(ctx, input.Owner, input.Repo, input.PRNumber, &clientv1.CreatePRCommentInput{
		Body:     input.Body,
		CommitID: input.CommitID,
		Path:     input.Path,
		Line:     input.Line,
	})
	if err != nil {
		return nil, fmt.Errorf("creating line comment: %w", err)
	}
	return &CommentResult{
		ID:      comment.ID,
		HTMLURL: comment.HTMLURL,
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
	ghPR, err := c.client.GetPullRequest(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}
	return &PRInfo{
		Number:  ghPR.Number,
		Title:   ghPR.Title,
		Body:    ghPR.Body,
		State:   ghPR.State,
		Author:  userLogin(ghPR.User),
		Head:    branchRef(ghPR.Head),
		Base:    branchRef(ghPR.Base),
		Commits: ghPR.Commits,
		HTMLURL: ghPR.HTMLURL,
	}, nil
}

// GetPRDiff retrieves the diff for a pull request.
func (c *Client) GetPRDiff(ctx context.Context, owner, repo string, number int) (string, error) {
	return c.client.GetPullRequestDiff(ctx, owner, repo, number)
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
	prs, err := c.client.ListPullRequests(ctx, owner, repo, &clientv1.ListPullRequestsOptions{
		State: "open",
	})
	if err != nil {
		return nil, err
	}

	result := make([]PRSummary, len(prs))
	for i, p := range prs {
		result[i] = PRSummary{
			Number:  p.Number,
			Title:   p.Title,
			Author:  userLogin(p.User),
			Head:    branchRef(p.Head),
			HTMLURL: p.HTMLURL,
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

func userLogin(u *gogithub.User) string {
	if u == nil {
		return ""
	}
	return u.Login
}

func branchRef(b *gogithub.PullRequestBranch) string {
	if b == nil {
		return ""
	}
	return b.Ref
}
