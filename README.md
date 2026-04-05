# Agent Code Review

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/agent-code-review/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/agent-code-review/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/agent-code-review/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/agent-code-review/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/agent-code-review/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/agent-code-review/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/agent-code-review
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/agent-code-review
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/agent-code-review
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/agent-code-review
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fagent-code-review
 [loc-svg]: https://tokei.rs/b1/github/plexusone/agent-code-review
 [repo-url]: https://github.com/plexusone/agent-code-review
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/agent-code-review/blob/master/LICENSE

AI-powered code review agent for GitHub Pull Requests.

[![Powered by Claude](https://img.shields.io/badge/Powered%20by-Claude-blueviolet)](https://anthropic.com/claude)

## Overview

Agent Code Review is a tool for posting AI-assisted code reviews to GitHub PRs. It supports three usage modes:

1. **Go SDK** — Programmatic access for integration into other tools
2. **CLI** — Command-line interface for scripts and manual use
3. **MCP Server** — Model Context Protocol server for Claude Code integration

Reviews are posted as a GitHub App, appearing as `PlexusOne Code Review[bot]`, clearly distinguishing AI-assisted reviews from human reviews.

> **Note:** This is a community project by PlexusOne. It is not an official Anthropic product.

## Features

- 🤖 **GitHub App Integration** — Reviews posted as a bot, not your personal account
- 📦 **Go SDK** — Programmatic API via `pkg/review`
- ⌨️ **CLI** — Full-featured command-line interface
- 🔌 **MCP Server** — Native integration with Claude Code
- 📋 **Multi-Agent Spec** — Agent definition follows the [multi-agent-spec](https://github.com/plexusone/multi-agent-spec) format

## Installation

### CLI

```bash
go install github.com/plexusone/agent-code-review/cmd/acr@latest
```

### Go SDK

```bash
go get github.com/plexusone/agent-code-review
```

## Prerequisites

### Option 1: GitHub App (Recommended)

1. Create a GitHub App (see [GitHub App Setup](#github-app-setup))
2. Configure authentication:

   ```bash
   export GITHUB_APP_ID=123456
   export GITHUB_INSTALLATION_ID=12345678
   export GITHUB_PRIVATE_KEY_PATH=~/.config/gogithub/private-key.pem
   ```

   Or create `~/.config/gogithub/app.json`:

   ```json
   {
     "app_id": 123456,
     "installation_id": 12345678,
     "private_key_path": "~/.config/gogithub/private-key.pem"
   }
   ```

### Option 2: Personal Access Token

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
```

## CLI Usage

```bash
# List open PRs
acr list -o owner -r repo

# Get PR details
acr get 123 -o owner -r repo

# Get PR diff
acr diff 123 -o owner -r repo

# Post a review
acr review 123 -o owner -r repo -e APPROVE -b "LGTM!"
acr review 123 -o owner -r repo -e REQUEST_CHANGES -f review.md
echo "Great work!" | acr review 123 -o owner -r repo -e COMMENT

# Add a comment
acr comment 123 -o owner -r repo -b "Thanks for the contribution!"

# Run as MCP server
acr serve
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `GITHUB_OWNER` | Default repository owner |
| `GITHUB_REPO` | Default repository name |
| `GITHUB_TOKEN` | Personal access token (fallback auth) |
| `GITHUB_APP_ID` | GitHub App ID |
| `GITHUB_INSTALLATION_ID` | GitHub App installation ID |
| `GITHUB_PRIVATE_KEY_PATH` | Path to GitHub App private key |

## Go SDK Usage

```go
package main

import (
    "context"
    "log"

    "github.com/grokify/gogithub/auth"
    "github.com/plexusone/agent-code-review/pkg/review"
)

func main() {
    ctx := context.Background()

    // Create client with GitHub App auth
    cfg, _ := auth.LoadAppConfig()
    client, _ := review.NewClientFromAppConfig(ctx, cfg)

    // Or with token auth
    // client := review.NewClientFromToken(ctx, "ghp_xxxx")

    // Post a review
    result, err := client.CreateReview(ctx, &review.ReviewInput{
        Owner:    "owner",
        Repo:     "repo",
        PRNumber: 123,
        Event:    review.EventApprove,
        Body:     "LGTM! Great work.",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Review posted: %s", result.HTMLURL)
}
```

## MCP Server Usage

Add to your Claude Code MCP configuration:

```json
{
  "mcpServers": {
    "code-review": {
      "command": "acr",
      "args": ["serve"]
    }
  }
}
```

### MCP Tools

| Tool | Description |
|------|-------------|
| `review_pr` | Post a full code review (approve/comment/request-changes) |
| `comment_pr` | Add a general comment to a PR |
| `line_comment` | Comment on a specific line in the diff |
| `get_pr_diff` | Fetch PR diff for analysis |
| `get_pr` | Get PR metadata |
| `list_prs` | List open PRs in a repository |

## GitHub App Setup

1. Go to **Settings → Developer settings → GitHub Apps → New GitHub App**

2. Configure:

   | Field | Value |
   |-------|-------|
   | Name | `PlexusOne Code Review` (or your preferred name) |
   | Homepage URL | `https://github.com/plexusone/agent-code-review` |
   | Webhook | Disable (not needed for local use) |

3. Set permissions:

   | Permission | Access |
   |------------|--------|
   | Pull requests | Read & Write |
   | Contents | Read |
   | Metadata | Read |

4. Generate and download a **private key**

5. Install the app on your repositories

6. Note your **App ID** and **Installation ID**

## Review Output Format

Reviews include a footer for transparency:

```markdown
## Code Review Summary

[Review content...]

---
🤖 Powered by Claude • PlexusOne Code Review
```

## Project Structure

```
agent-code-review/
├── cmd/acr/                 # CLI application
│   ├── main.go
│   └── cmd/
│       ├── root.go
│       ├── review.go
│       ├── comment.go
│       ├── diff.go
│       ├── get.go
│       ├── list.go
│       └── serve.go
├── pkg/review/              # Go SDK
│   └── review.go
├── internal/mcp/            # MCP server implementation
│   └── server.go
├── specs/agents/            # Agent definition (multi-agent-spec)
│   └── code-reviewer.md
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Dependencies

- [gogithub](https://github.com/grokify/gogithub) — GitHub API wrapper with App authentication
- [go-sdk](https://github.com/modelcontextprotocol/go-sdk) — Official MCP Go SDK
- [cobra](https://github.com/spf13/cobra) — CLI framework

## Related Projects

- [multi-agent-spec](https://github.com/plexusone/multi-agent-spec) — Portable multi-agent specification format
- [agent-a11y](https://github.com/plexusone/agent-a11y) — Accessibility testing agent
- [agent-dast](https://github.com/plexusone/agent-dast) — Dynamic application security testing agent
- [gogithub](https://github.com/grokify/gogithub) — GitHub API wrapper SDK

## License

MIT

---

<sub>🤖 Powered by [Claude](https://anthropic.com/claude) — Built by [PlexusOne](https://github.com/plexusone)</sub>
