# CLI Reference

The `acr` command-line interface provides full access to code review operations.

## Global Flags

These flags are available on all commands:

| Flag | Short | Description |
|------|-------|-------------|
| `--owner` | `-o` | Repository owner (user or organization) |
| `--repo` | `-r` | Repository name |

You can also set defaults via environment variables:

```bash
export GITHUB_OWNER=myorg
export GITHUB_REPO=myrepo
```

## Commands

### list

List open pull requests in a repository.

```bash
acr list -o owner -r repo
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON |

**Examples:**

```bash
# List PRs in table format
acr list -o plexusone -r agent-code-review

# Output as JSON for scripting
acr list -o plexusone -r agent-code-review --json
```

---

### get

Get details for a specific pull request.

```bash
acr get <pr-number> -o owner -r repo
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON |

**Examples:**

```bash
# Get PR details
acr get 123 -o plexusone -r agent-code-review

# Get as JSON
acr get 123 -o plexusone -r agent-code-review --json
```

**Output:**

```
PR #123: Add new feature
Author: contributor
State: open
Branch: feature-branch -> main
Commits: 3
URL: https://github.com/plexusone/agent-code-review/pull/123

This PR adds a new feature that...
```

---

### diff

Fetch the diff for a pull request.

```bash
acr diff <pr-number> -o owner -r repo
```

**Examples:**

```bash
# Print diff to stdout
acr diff 123 -o plexusone -r agent-code-review

# Save to file
acr diff 123 -o plexusone -r agent-code-review > pr-123.diff

# Pipe to another tool
acr diff 123 -o plexusone -r agent-code-review | less
```

---

### review

Post a code review to a pull request.

```bash
acr review <pr-number> -o owner -r repo -e <event> -b <body>
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--event` | `-e` | Review event: `APPROVE`, `REQUEST_CHANGES`, or `COMMENT` |
| `--body` | `-b` | Review body text |
| `--file` | `-f` | Read review body from file |

**Input Sources:**

The review body can be provided in three ways (in order of precedence):

1. `--body` flag — Direct text
2. `--file` flag — Read from file
3. `stdin` — Pipe content

**Examples:**

```bash
# Approve with inline body
acr review 123 -o owner -r repo -e APPROVE -b "LGTM! Great work."

# Request changes with body from file
acr review 123 -o owner -r repo -e REQUEST_CHANGES -f review.md

# Comment with body from stdin
echo "Thanks for the contribution!" | acr review 123 -o owner -r repo -e COMMENT

# Pipe from another command (e.g., AI-generated review)
cat review-output.md | acr review 123 -o owner -r repo -e APPROVE
```

**Review Events:**

| Event | Description |
|-------|-------------|
| `APPROVE` | Approve the pull request |
| `REQUEST_CHANGES` | Request changes before merging |
| `COMMENT` | Add a review comment without approval/rejection |

---

### comment

Add a general comment to a pull request (not a review).

```bash
acr comment <pr-number> -o owner -r repo -b <body>
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--body` | `-b` | Comment body text |
| `--file` | `-f` | Read comment body from file |

**Examples:**

```bash
# Add a comment
acr comment 123 -o owner -r repo -b "Thanks for the contribution!"

# Comment from file
acr comment 123 -o owner -r repo -f comment.md
```

!!! note
    This creates an issue comment, not a review. Use `acr review` for formal code reviews.

---

### serve

Run as an MCP (Model Context Protocol) server.

```bash
acr serve
```

The server communicates over stdin/stdout using JSON-RPC. See [MCP Integration](mcp.md) for configuration details.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (authentication, API, validation) |

## Scripting Examples

### Review All Open PRs

```bash
#!/bin/bash
for pr in $(acr list -o owner -r repo --json | jq -r '.[].number'); do
    echo "Reviewing PR #$pr..."
    acr diff $pr -o owner -r repo > /tmp/pr-$pr.diff
    # Process diff and generate review...
done
```

### Conditional Approval

```bash
#!/bin/bash
PR=$1
DIFF=$(acr diff $PR -o owner -r repo)

if echo "$DIFF" | grep -q "TODO"; then
    acr review $PR -o owner -r repo -e REQUEST_CHANGES -b "Please resolve TODOs before merging."
else
    acr review $PR -o owner -r repo -e APPROVE -b "LGTM!"
fi
```
