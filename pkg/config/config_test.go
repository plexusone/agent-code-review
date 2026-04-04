package config

import (
	"os"
	"testing"
)

func TestResolveRepo_ExplicitValues(t *testing.T) {
	cfg, err := ResolveRepo("myowner", "myrepo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Owner != "myowner" {
		t.Errorf("Owner = %s, want myowner", cfg.Owner)
	}
	if cfg.Repo != "myrepo" {
		t.Errorf("Repo = %s, want myrepo", cfg.Repo)
	}
}

func TestResolveRepo_FromEnv(t *testing.T) {
	// Save and restore env vars
	oldOwner := os.Getenv("GITHUB_OWNER")
	oldRepo := os.Getenv("GITHUB_REPO")
	defer func() {
		os.Setenv("GITHUB_OWNER", oldOwner)
		os.Setenv("GITHUB_REPO", oldRepo)
	}()

	os.Setenv("GITHUB_OWNER", "envowner")
	os.Setenv("GITHUB_REPO", "envrepo")

	cfg, err := ResolveRepo("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Owner != "envowner" {
		t.Errorf("Owner = %s, want envowner", cfg.Owner)
	}
	if cfg.Repo != "envrepo" {
		t.Errorf("Repo = %s, want envrepo", cfg.Repo)
	}
}

func TestResolveRepo_ExplicitOverridesEnv(t *testing.T) {
	// Save and restore env vars
	oldOwner := os.Getenv("GITHUB_OWNER")
	oldRepo := os.Getenv("GITHUB_REPO")
	defer func() {
		os.Setenv("GITHUB_OWNER", oldOwner)
		os.Setenv("GITHUB_REPO", oldRepo)
	}()

	os.Setenv("GITHUB_OWNER", "envowner")
	os.Setenv("GITHUB_REPO", "envrepo")

	cfg, err := ResolveRepo("explicit", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Owner != "explicit" {
		t.Errorf("Owner = %s, want explicit", cfg.Owner)
	}
	if cfg.Repo != "envrepo" {
		t.Errorf("Repo = %s, want envrepo", cfg.Repo)
	}
}

func TestResolveRepo_MissingOwner(t *testing.T) {
	// Save and restore env vars
	oldOwner := os.Getenv("GITHUB_OWNER")
	oldRepo := os.Getenv("GITHUB_REPO")
	defer func() {
		os.Setenv("GITHUB_OWNER", oldOwner)
		os.Setenv("GITHUB_REPO", oldRepo)
	}()

	os.Unsetenv("GITHUB_OWNER")
	os.Unsetenv("GITHUB_REPO")

	_, err := ResolveRepo("", "myrepo")
	if err == nil {
		t.Fatal("expected error for missing owner")
	}
}

func TestResolveRepo_MissingRepo(t *testing.T) {
	// Save and restore env vars
	oldOwner := os.Getenv("GITHUB_OWNER")
	oldRepo := os.Getenv("GITHUB_REPO")
	defer func() {
		os.Setenv("GITHUB_OWNER", oldOwner)
		os.Setenv("GITHUB_REPO", oldRepo)
	}()

	os.Unsetenv("GITHUB_OWNER")
	os.Unsetenv("GITHUB_REPO")

	_, err := ResolveRepo("myowner", "")
	if err == nil {
		t.Fatal("expected error for missing repo")
	}
}

func TestResolveRepo_BothMissing(t *testing.T) {
	// Save and restore env vars
	oldOwner := os.Getenv("GITHUB_OWNER")
	oldRepo := os.Getenv("GITHUB_REPO")
	defer func() {
		os.Setenv("GITHUB_OWNER", oldOwner)
		os.Setenv("GITHUB_REPO", oldRepo)
	}()

	os.Unsetenv("GITHUB_OWNER")
	os.Unsetenv("GITHUB_REPO")

	_, err := ResolveRepo("", "")
	if err == nil {
		t.Fatal("expected error for missing owner and repo")
	}
}
