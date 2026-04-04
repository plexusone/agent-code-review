# Go SDK

The `pkg/review` package provides a high-level Go API for GitHub code review operations.

## Installation

```bash
go get github.com/plexusone/agent-code-review
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/plexusone/agent-code-review/pkg/review"
)

func main() {
    ctx := context.Background()

    // Create client with token auth
    client := review.NewClientFromToken(ctx, "ghp_xxxxxxxxxxxx")

    // Approve a PR
    result, err := client.Approve(ctx, "owner", "repo", 123, "LGTM!")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Review posted: %s", result.HTMLURL)
}
```

## Creating a Client

### From Personal Access Token

```go
client := review.NewClientFromToken(ctx, token)
```

### From GitHub App Config

```go
import "github.com/grokify/gogithub/auth"

cfg, err := auth.LoadAppConfig()
if err != nil {
    log.Fatal(err)
}

client, err := review.NewClientFromAppConfig(ctx, cfg)
if err != nil {
    log.Fatal(err)
}
```

### From Existing GitHub Client

```go
import "github.com/google/go-github/v84/github"

gh := github.NewClient(nil).WithAuthToken(token)
client := review.NewClient(gh)
```

## Review Operations

### CreateReview

Post a full code review to a pull request.

```go
result, err := client.CreateReview(ctx, &review.ReviewInput{
    Owner:    "owner",
    Repo:     "repo",
    PRNumber: 123,
    Event:    review.EventApprove, // or EventRequestChanges, EventComment
    Body:     "Great work! The implementation looks solid.",
})
```

### Convenience Methods

```go
// Approve a PR
result, err := client.Approve(ctx, "owner", "repo", 123, "LGTM!")

// Request changes
result, err := client.RequestChanges(ctx, "owner", "repo", 123, "Please fix the security issue.")

// Add a review comment (no approval/rejection)
result, err := client.Comment(ctx, "owner", "repo", 123, "Some observations...")
```

### Review Events

| Constant | Value | Description |
|----------|-------|-------------|
| `EventApprove` | `"APPROVE"` | Approve the pull request |
| `EventRequestChanges` | `"REQUEST_CHANGES"` | Request changes |
| `EventComment` | `"COMMENT"` | Comment without verdict |

## Comment Operations

### CreateComment

Add a general comment to a PR (issue comment, not a review).

```go
result, err := client.CreateComment(ctx, &review.CommentInput{
    Owner:    "owner",
    Repo:     "repo",
    PRNumber: 123,
    Body:     "Thanks for the contribution!",
})
```

### CreateLineComment

Add a comment on a specific line in the diff.

```go
result, err := client.CreateLineComment(ctx, &review.LineCommentInput{
    Owner:    "owner",
    Repo:     "repo",
    PRNumber: 123,
    CommitID: "abc123def456...",
    Path:     "src/main.go",
    Line:     42,
    Body:     "Consider using a more descriptive variable name here.",
})
```

## Query Operations

### GetPR

Retrieve pull request details.

```go
pr, err := client.GetPR(ctx, "owner", "repo", 123)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("PR #%d: %s\n", pr.Number, pr.Title)
fmt.Printf("Author: %s\n", pr.Author)
fmt.Printf("State: %s\n", pr.State)
fmt.Printf("Branch: %s -> %s\n", pr.Head, pr.Base)
```

**PRInfo Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `Number` | `int` | PR number |
| `Title` | `string` | PR title |
| `Body` | `string` | PR description |
| `State` | `string` | `"open"` or `"closed"` |
| `Author` | `string` | Author's GitHub username |
| `Head` | `string` | Source branch |
| `Base` | `string` | Target branch |
| `Commits` | `int` | Number of commits |
| `HTMLURL` | `string` | Web URL |

### GetPRDiff

Retrieve the diff for a pull request.

```go
diff, err := client.GetPRDiff(ctx, "owner", "repo", 123)
if err != nil {
    log.Fatal(err)
}

fmt.Println(diff)
```

### ListOpenPRs

List open pull requests in a repository.

```go
prs, err := client.ListOpenPRs(ctx, "owner", "repo")
if err != nil {
    log.Fatal(err)
}

for _, pr := range prs {
    fmt.Printf("#%d %s (%s)\n", pr.Number, pr.Title, pr.Author)
}
```

## Review Footer

All reviews automatically include a footer for transparency:

```markdown
---
🤖 Powered by Claude • PlexusOne Code Review
```

This is defined as `review.ReviewFooter` and appended to all review bodies.

## Error Handling

All operations return wrapped errors with context:

```go
result, err := client.CreateReview(ctx, input)
if err != nil {
    // Error includes context: "creating review: POST .../reviews: 404 Not Found"
    log.Printf("Failed to create review: %v", err)
    return err
}
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/grokify/gogithub/auth"
    "github.com/plexusone/agent-code-review/pkg/review"
)

func main() {
    ctx := context.Background()

    // Create client
    cfg, _ := auth.LoadAppConfig()
    client, err := review.NewClientFromAppConfig(ctx, cfg)
    if err != nil {
        log.Fatal(err)
    }

    // Get PR diff
    diff, err := client.GetPRDiff(ctx, "owner", "repo", 123)
    if err != nil {
        log.Fatal(err)
    }

    // Analyze diff (simplified example)
    var event review.ReviewEvent
    var body string

    if strings.Contains(diff, "TODO") {
        event = review.EventRequestChanges
        body = "Please resolve TODOs before merging."
    } else {
        event = review.EventApprove
        body = "LGTM! No issues found."
    }

    // Post review
    result, err := client.CreateReview(ctx, &review.ReviewInput{
        Owner:    "owner",
        Repo:     "repo",
        PRNumber: 123,
        Event:    event,
        Body:     body,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Review posted: %s\n", result.HTMLURL)
}
```
