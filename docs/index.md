# Agent Code Review

AI-powered code review agent for GitHub Pull Requests.

[![Powered by Claude](https://img.shields.io/badge/Powered%20by-Claude-blueviolet)](https://anthropic.com/claude)

## Overview

Agent Code Review is a tool for posting AI-assisted code reviews to GitHub PRs. It supports three usage modes:

1. **Go SDK** — Programmatic access for integration into other tools
2. **CLI** — Command-line interface for scripts and manual use
3. **MCP Server** — Model Context Protocol server for Claude Code integration

Reviews are posted as a GitHub App, appearing as `PlexusOne Code Review[bot]`, clearly distinguishing AI-assisted reviews from human reviews.

!!! note
    This is a community project by PlexusOne. It is not an official Anthropic product.

## Features

- **GitHub App Integration** — Reviews posted as a bot, not your personal account
- **Go SDK** — Programmatic API via `pkg/review`
- **CLI** — Full-featured command-line interface
- **MCP Server** — Native integration with Claude Code
- **Multi-Agent Spec** — Agent definition follows the [multi-agent-spec](https://github.com/plexusone/multi-agent-spec) format

## Quick Start

=== "CLI"

    ```bash
    # Install
    go install github.com/plexusone/agent-code-review/cmd/acr@latest

    # Set up auth
    export GITHUB_TOKEN=ghp_xxxxxxxxxxxx

    # Review a PR
    acr review 123 -o owner -r repo -e APPROVE -b "LGTM!"
    ```

=== "Go SDK"

    ```go
    import "github.com/plexusone/agent-code-review/pkg/review"

    client := review.NewClientFromToken(ctx, token)
    result, _ := client.Approve(ctx, "owner", "repo", 123, "LGTM!")
    ```

=== "MCP Server"

    Add to Claude Code MCP config:

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

## Project Structure

```
agent-code-review/
├── cmd/acr/                 # CLI application
├── pkg/review/              # Go SDK
├── pkg/config/              # Configuration utilities
├── pkg/input/               # Input parsing utilities
├── internal/mcp/            # MCP server implementation
├── specs/agents/            # Agent definition (multi-agent-spec)
└── docs/                    # Documentation (this site)
```

## Related Projects

- [multi-agent-spec](https://github.com/plexusone/multi-agent-spec) — Portable multi-agent specification format
- [gogithub](https://github.com/grokify/gogithub) — GitHub API wrapper SDK
- [go-sdk](https://github.com/modelcontextprotocol/go-sdk) — Official MCP Go SDK

## License

MIT
