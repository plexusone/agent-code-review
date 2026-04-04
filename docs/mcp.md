# MCP Integration

Agent Code Review includes an MCP (Model Context Protocol) server for native integration with Claude Code and other MCP-compatible clients.

## Overview

The MCP server exposes code review functionality as tools that Claude can use directly within Claude Code. This enables AI-assisted code reviews without leaving your development environment.

## Configuration

### Claude Code

Add the following to your Claude Code MCP configuration:

=== "macOS/Linux"

    Edit `~/.config/claude-code/mcp.json`:

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

=== "Windows"

    Edit `%APPDATA%\claude-code\mcp.json`:

    ```json
    {
      "mcpServers": {
        "code-review": {
          "command": "acr.exe",
          "args": ["serve"]
        }
      }
    }
    ```

### With Environment Variables

If you need to pass authentication via environment variables:

```json
{
  "mcpServers": {
    "code-review": {
      "command": "acr",
      "args": ["serve"],
      "env": {
        "GITHUB_TOKEN": "ghp_xxxxxxxxxxxx"
      }
    }
  }
}
```

!!! tip
    For GitHub App authentication, ensure the config file exists at `~/.config/gogithub/app.json` or set the environment variables before starting Claude Code.

## Available Tools

The MCP server exposes six tools:

### review_pr

Post a code review to a GitHub pull request.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `owner` | string | Yes | Repository owner |
| `repo` | string | Yes | Repository name |
| `pr_number` | integer | Yes | Pull request number |
| `event` | string | Yes | `APPROVE`, `REQUEST_CHANGES`, or `COMMENT` |
| `body` | string | Yes | Review body (markdown) |

**Example Usage in Claude:**

> "Review PR #123 in plexusone/agent-code-review and approve it with the message 'LGTM!'"

---

### comment_pr

Add a general comment to a pull request.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `owner` | string | Yes | Repository owner |
| `repo` | string | Yes | Repository name |
| `pr_number` | integer | Yes | Pull request number |
| `body` | string | Yes | Comment body (markdown) |

**Example Usage:**

> "Add a comment to PR #123 thanking the contributor"

---

### line_comment

Add a comment on a specific line in the PR diff.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `owner` | string | Yes | Repository owner |
| `repo` | string | Yes | Repository name |
| `pr_number` | integer | Yes | Pull request number |
| `commit_id` | string | Yes | Commit SHA to comment on |
| `path` | string | Yes | File path relative to repo root |
| `line` | integer | Yes | Line number in the diff |
| `body` | string | Yes | Comment body (markdown) |

**Example Usage:**

> "Add a line comment on main.go line 42 suggesting a better variable name"

---

### get_pr_diff

Fetch the diff for a pull request.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `owner` | string | Yes | Repository owner |
| `repo` | string | Yes | Repository name |
| `pr_number` | integer | Yes | Pull request number |

**Returns:** The unified diff as text.

**Example Usage:**

> "Get the diff for PR #123 so I can review it"

---

### get_pr

Get pull request metadata.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `owner` | string | Yes | Repository owner |
| `repo` | string | Yes | Repository name |
| `pr_number` | integer | Yes | Pull request number |

**Returns:** PR details including title, author, state, branches, and commit count.

**Example Usage:**

> "What's the status of PR #123?"

---

### list_prs

List open pull requests in a repository.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `owner` | string | Yes | Repository owner |
| `repo` | string | Yes | Repository name |

**Returns:** List of open PRs with number, title, author, and branch.

**Example Usage:**

> "List the open PRs in plexusone/agent-code-review"

## Workflow Example

A typical code review workflow in Claude Code:

1. **List PRs** — "Show me open PRs in my-org/my-repo"
2. **Get Diff** — "Get the diff for PR #42"
3. **Analyze** — Claude reviews the code changes
4. **Post Review** — "Approve PR #42 with a summary of your findings"

## Troubleshooting

### Server Not Starting

Check that `acr` is in your PATH:

```bash
which acr
```

If not found, ensure the Go bin directory is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Authentication Errors

Verify authentication is configured:

```bash
# Test with CLI
acr list -o owner -r repo
```

If this works, the MCP server will use the same authentication.

### Debug Mode

Run the server manually to see output:

```bash
acr serve
```

The server communicates via stdin/stdout using JSON-RPC. Any errors will be returned in the response.
