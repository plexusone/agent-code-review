# Getting Started

This guide covers installation and authentication setup for Agent Code Review.

## Installation

### CLI

Install the `acr` command-line tool:

```bash
go install github.com/plexusone/agent-code-review/cmd/acr@latest
```

Verify installation:

```bash
acr --help
```

### Go SDK

Add the SDK to your Go project:

```bash
go get github.com/plexusone/agent-code-review
```

## Authentication

Agent Code Review supports two authentication methods. Choose the one that fits your use case.

### Option 1: GitHub App (Recommended)

Using a GitHub App provides several benefits:

- Reviews appear as `PlexusOne Code Review[bot]` (or your app name)
- Clear distinction between AI and human reviews
- Fine-grained permissions per repository
- No personal account rate limits

#### Creating a GitHub App

1. Go to **Settings → Developer settings → GitHub Apps → New GitHub App**

2. Configure the app:

    | Field | Value |
    |-------|-------|
    | Name | `PlexusOne Code Review` (or your preferred name) |
    | Homepage URL | `https://github.com/plexusone/agent-code-review` |
    | Webhook | Disable (uncheck "Active") |

3. Set permissions:

    | Permission | Access |
    |------------|--------|
    | Pull requests | Read & Write |
    | Contents | Read |
    | Metadata | Read |

4. Click **Create GitHub App**

5. Generate a **private key** and download it

6. **Install** the app on your repositories

7. Note your **App ID** (shown on the app settings page) and **Installation ID** (from the URL after installing)

#### Configuring Authentication

=== "Environment Variables"

    ```bash
    export GITHUB_APP_ID=123456
    export GITHUB_INSTALLATION_ID=12345678
    export GITHUB_PRIVATE_KEY_PATH=~/.config/gogithub/private-key.pem
    ```

=== "Config File"

    Create `~/.config/gogithub/app.json`:

    ```json
    {
      "app_id": 123456,
      "installation_id": 12345678,
      "private_key_path": "~/.config/gogithub/private-key.pem"
    }
    ```

### Option 2: Personal Access Token

For quick setup or personal use, you can use a personal access token:

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
```

!!! warning
    Reviews will appear as coming from your personal account, not a bot.

#### Creating a Token

1. Go to **Settings → Developer settings → Personal access tokens → Tokens (classic)**
2. Click **Generate new token (classic)**
3. Select scopes:
    - `repo` (Full control of private repositories)
4. Generate and copy the token

## Environment Variables Reference

| Variable | Description |
|----------|-------------|
| `GITHUB_TOKEN` | Personal access token (fallback auth) |
| `GITHUB_APP_ID` | GitHub App ID |
| `GITHUB_INSTALLATION_ID` | GitHub App installation ID |
| `GITHUB_PRIVATE_KEY_PATH` | Path to GitHub App private key |
| `GITHUB_OWNER` | Default repository owner |
| `GITHUB_REPO` | Default repository name |

## Verifying Setup

Test your authentication by listing PRs:

```bash
acr list -o owner -r repo
```

If configured correctly, you'll see a list of open pull requests (or a message indicating no open PRs).

## Next Steps

- [CLI Reference](cli.md) — Learn the available commands
- [Go SDK](sdk.md) — Integrate into your Go applications
- [MCP Integration](mcp.md) — Use with Claude Code
