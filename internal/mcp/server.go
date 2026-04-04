// Package mcp provides the Model Context Protocol server implementation.
package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plexusone/agent-code-review/pkg/review"
)

// Server implements the MCP server for code review tools.
type Server struct {
	client    *review.Client
	mcpServer *mcp.Server
}

// NewServer creates a new MCP server with the given review client.
func NewServer(client *review.Client) *Server {
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "agent-code-review",
			Version: "0.1.0",
		},
		nil,
	)

	s := &Server{
		client:    client,
		mcpServer: mcpServer,
	}

	s.registerTools()

	return s
}

// Run starts the MCP server on stdin/stdout.
func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

// Tool input/output types with jsonschema annotations

// ReviewPRInput is the input for the review_pr tool.
type ReviewPRInput struct {
	Owner    string `json:"owner" jsonschema:"description=Repository owner (user or org)"`
	Repo     string `json:"repo" jsonschema:"description=Repository name"`
	PRNumber int    `json:"pr_number" jsonschema:"description=Pull request number"`
	Event    string `json:"event" jsonschema:"description=Review action,enum=APPROVE,enum=REQUEST_CHANGES,enum=COMMENT"`
	Body     string `json:"body" jsonschema:"description=Review body (markdown)"`
}

// ReviewPROutput is the output for the review_pr tool.
type ReviewPROutput struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}

// CommentPRInput is the input for the comment_pr tool.
type CommentPRInput struct {
	Owner    string `json:"owner" jsonschema:"description=Repository owner (user or org)"`
	Repo     string `json:"repo" jsonschema:"description=Repository name"`
	PRNumber int    `json:"pr_number" jsonschema:"description=Pull request number"`
	Body     string `json:"body" jsonschema:"description=Comment body (markdown)"`
}

// CommentPROutput is the output for the comment_pr tool.
type CommentPROutput struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}

// LineCommentInput is the input for the line_comment tool.
type LineCommentInput struct {
	Owner    string `json:"owner" jsonschema:"description=Repository owner (user or org)"`
	Repo     string `json:"repo" jsonschema:"description=Repository name"`
	PRNumber int    `json:"pr_number" jsonschema:"description=Pull request number"`
	CommitID string `json:"commit_id" jsonschema:"description=Commit SHA to comment on"`
	Path     string `json:"path" jsonschema:"description=File path relative to repo root"`
	Line     int    `json:"line" jsonschema:"description=Line number in the diff"`
	Body     string `json:"body" jsonschema:"description=Comment body (markdown)"`
}

// LineCommentOutput is the output for the line_comment tool.
type LineCommentOutput struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}

// GetPRDiffInput is the input for the get_pr_diff tool.
type GetPRDiffInput struct {
	Owner    string `json:"owner" jsonschema:"description=Repository owner (user or org)"`
	Repo     string `json:"repo" jsonschema:"description=Repository name"`
	PRNumber int    `json:"pr_number" jsonschema:"description=Pull request number"`
}

// GetPRDiffOutput is the output for the get_pr_diff tool.
type GetPRDiffOutput struct {
	Diff string `json:"diff"`
}

// GetPRInput is the input for the get_pr tool.
type GetPRInput struct {
	Owner    string `json:"owner" jsonschema:"description=Repository owner (user or org)"`
	Repo     string `json:"repo" jsonschema:"description=Repository name"`
	PRNumber int    `json:"pr_number" jsonschema:"description=Pull request number"`
}

// GetPROutput is the output for the get_pr tool.
type GetPROutput struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	State   string `json:"state"`
	Author  string `json:"author"`
	Head    string `json:"head"`
	Base    string `json:"base"`
	Commits int    `json:"commits"`
	URL     string `json:"url"`
}

// ListPRsInput is the input for the list_prs tool.
type ListPRsInput struct {
	Owner string `json:"owner" jsonschema:"description=Repository owner (user or org)"`
	Repo  string `json:"repo" jsonschema:"description=Repository name"`
}

// PRSummary is a summary of a PR for the list_prs output.
type PRSummary struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Head   string `json:"head"`
	URL    string `json:"url"`
}

// ListPRsOutput is the output for the list_prs tool.
type ListPRsOutput struct {
	PullRequests []PRSummary `json:"pull_requests"`
}

func (s *Server) registerTools() {
	// review_pr tool
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{
			Name:        "review_pr",
			Description: "Post a code review to a GitHub pull request",
		},
		s.reviewPR,
	)

	// comment_pr tool
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{
			Name:        "comment_pr",
			Description: "Add a general comment to a pull request",
		},
		s.commentPR,
	)

	// line_comment tool
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{
			Name:        "line_comment",
			Description: "Add a comment on a specific line in a PR diff",
		},
		s.lineComment,
	)

	// get_pr_diff tool
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{
			Name:        "get_pr_diff",
			Description: "Fetch the diff for a pull request",
		},
		s.getPRDiff,
	)

	// get_pr tool
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{
			Name:        "get_pr",
			Description: "Get pull request details (title, body, author, files)",
		},
		s.getPR,
	)

	// list_prs tool
	mcp.AddTool(s.mcpServer,
		&mcp.Tool{
			Name:        "list_prs",
			Description: "List open pull requests in a repository",
		},
		s.listPRs,
	)
}

func (s *Server) reviewPR(ctx context.Context, req *mcp.CallToolRequest, input ReviewPRInput) (*mcp.CallToolResult, ReviewPROutput, error) {
	result, err := s.client.CreateReview(ctx, &review.ReviewInput{
		Owner:    input.Owner,
		Repo:     input.Repo,
		PRNumber: input.PRNumber,
		Event:    review.ReviewEvent(input.Event),
		Body:     input.Body,
	})
	if err != nil {
		return nil, ReviewPROutput{}, err
	}

	output := ReviewPROutput{
		Message: fmt.Sprintf("Review posted to %s/%s#%d with event %s", input.Owner, input.Repo, input.PRNumber, input.Event),
		URL:     result.HTMLURL,
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output.Message},
		},
	}, output, nil
}

func (s *Server) commentPR(ctx context.Context, req *mcp.CallToolRequest, input CommentPRInput) (*mcp.CallToolResult, CommentPROutput, error) {
	result, err := s.client.CreateComment(ctx, &review.CommentInput{
		Owner:    input.Owner,
		Repo:     input.Repo,
		PRNumber: input.PRNumber,
		Body:     input.Body,
	})
	if err != nil {
		return nil, CommentPROutput{}, err
	}

	output := CommentPROutput{
		Message: fmt.Sprintf("Comment posted to %s/%s#%d", input.Owner, input.Repo, input.PRNumber),
		URL:     result.HTMLURL,
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output.Message},
		},
	}, output, nil
}

func (s *Server) lineComment(ctx context.Context, req *mcp.CallToolRequest, input LineCommentInput) (*mcp.CallToolResult, LineCommentOutput, error) {
	result, err := s.client.CreateLineComment(ctx, &review.LineCommentInput{
		Owner:    input.Owner,
		Repo:     input.Repo,
		PRNumber: input.PRNumber,
		CommitID: input.CommitID,
		Path:     input.Path,
		Line:     input.Line,
		Body:     input.Body,
	})
	if err != nil {
		return nil, LineCommentOutput{}, err
	}

	output := LineCommentOutput{
		Message: fmt.Sprintf("Line comment posted to %s/%s#%d at %s:%d", input.Owner, input.Repo, input.PRNumber, input.Path, input.Line),
		URL:     result.HTMLURL,
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output.Message},
		},
	}, output, nil
}

func (s *Server) getPRDiff(ctx context.Context, req *mcp.CallToolRequest, input GetPRDiffInput) (*mcp.CallToolResult, GetPRDiffOutput, error) {
	diff, err := s.client.GetPRDiff(ctx, input.Owner, input.Repo, input.PRNumber)
	if err != nil {
		return nil, GetPRDiffOutput{}, err
	}

	output := GetPRDiffOutput{Diff: diff}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: diff},
		},
	}, output, nil
}

func (s *Server) getPR(ctx context.Context, req *mcp.CallToolRequest, input GetPRInput) (*mcp.CallToolResult, GetPROutput, error) {
	pr, err := s.client.GetPR(ctx, input.Owner, input.Repo, input.PRNumber)
	if err != nil {
		return nil, GetPROutput{}, err
	}

	output := GetPROutput{
		Number:  pr.Number,
		Title:   pr.Title,
		Body:    pr.Body,
		State:   pr.State,
		Author:  pr.Author,
		Head:    pr.Head,
		Base:    pr.Base,
		Commits: pr.Commits,
		URL:     pr.HTMLURL,
	}

	text := fmt.Sprintf("PR #%d: %s\nAuthor: %s\nState: %s\nBranch: %s -> %s\nCommits: %d\nURL: %s",
		pr.Number, pr.Title, pr.Author, pr.State, pr.Head, pr.Base, pr.Commits, pr.HTMLURL)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}, output, nil
}

func (s *Server) listPRs(ctx context.Context, req *mcp.CallToolRequest, input ListPRsInput) (*mcp.CallToolResult, ListPRsOutput, error) {
	prs, err := s.client.ListOpenPRs(ctx, input.Owner, input.Repo)
	if err != nil {
		return nil, ListPRsOutput{}, err
	}

	output := ListPRsOutput{
		PullRequests: make([]PRSummary, len(prs)),
	}
	var text string
	for i, pr := range prs {
		output.PullRequests[i] = PRSummary{
			Number: pr.Number,
			Title:  pr.Title,
			Author: pr.Author,
			Head:   pr.Head,
			URL:    pr.HTMLURL,
		}
		text += fmt.Sprintf("#%-6d %s (%s)\n", pr.Number, pr.Title, pr.Author)
	}

	if len(prs) == 0 {
		text = "No open pull requests"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}, output, nil
}
