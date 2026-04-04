// Package config provides configuration resolution for agent-code-review.
package config

import (
	"context"
	"fmt"
	"os"

	"github.com/grokify/gogithub/auth"
	"github.com/plexusone/agent-code-review/pkg/review"
)

// RepoConfig holds repository identification.
type RepoConfig struct {
	Owner string
	Repo  string
}

// ResolveRepo resolves owner and repo from explicit values or environment variables.
// Explicit values take precedence over environment variables.
func ResolveRepo(owner, repo string) (*RepoConfig, error) {
	if owner == "" {
		owner = os.Getenv("GITHUB_OWNER")
	}
	if repo == "" {
		repo = os.Getenv("GITHUB_REPO")
	}

	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required (use flags or GITHUB_OWNER/GITHUB_REPO env vars)")
	}

	return &RepoConfig{Owner: owner, Repo: repo}, nil
}

// CreateClient creates a review client using available authentication.
// It tries GitHub App authentication first, then falls back to token authentication.
func CreateClient(ctx context.Context) (*review.Client, error) {
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
