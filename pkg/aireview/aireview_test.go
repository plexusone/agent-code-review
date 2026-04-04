package aireview

import (
	"strings"
	"testing"
)

func TestValidScopes(t *testing.T) {
	scopes := ValidScopes()
	if len(scopes) != 4 {
		t.Errorf("ValidScopes() returned %d scopes, want 4", len(scopes))
	}

	expected := map[Scope]bool{
		ScopeFull:        true,
		ScopeSecurity:    true,
		ScopeStyle:       true,
		ScopePerformance: true,
	}
	for _, s := range scopes {
		if !expected[s] {
			t.Errorf("unexpected scope: %s", s)
		}
	}
}

func TestIsValidScope(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"full", true},
		{"security", true},
		{"style", true},
		{"performance", true},
		{"FULL", false}, // case-sensitive
		{"invalid", false},
		{"", false},
	}
	for _, tt := range tests {
		got := IsValidScope(tt.input)
		if got != tt.want {
			t.Errorf("IsValidScope(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestBuildSystemPrompt(t *testing.T) {
	prompt := BuildSystemPrompt()
	if prompt == "" {
		t.Error("BuildSystemPrompt() returned empty string")
	}
	if !strings.Contains(prompt, "code reviewer") {
		t.Error("system prompt should mention code reviewer")
	}
}

func TestBuildUserPrompt(t *testing.T) {
	tests := []struct {
		scope   Scope
		title   string
		body    string
		diff    string
		wantIn  []string
		wantOut []string
	}{
		{
			scope:  ScopeFull,
			title:  "Add feature",
			body:   "This PR adds a new feature",
			diff:   "+func newFeature() {}",
			wantIn: []string{"Add feature", "This PR adds a new feature", "Correctness", "Security", "Performance"},
		},
		{
			scope:  ScopeSecurity,
			title:  "Fix auth",
			body:   "Security fix",
			diff:   "-password = 'secret'",
			wantIn: []string{"Fix auth", "Injection Vulnerabilities", "Authentication", "Sensitive Data"},
		},
		{
			scope:  ScopeStyle,
			title:  "Refactor",
			body:   "Cleanup",
			diff:   "+// better comment",
			wantIn: []string{"Refactor", "Naming", "Readability", "Consistency"},
		},
		{
			scope:  ScopePerformance,
			title:  "Optimize query",
			body:   "Perf improvement",
			diff:   "+WHERE indexed_col = ?",
			wantIn: []string{"Optimize query", "Database", "Algorithms", "Memory"},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.scope), func(t *testing.T) {
			prompt := BuildUserPrompt(tt.scope, tt.title, tt.body, tt.diff)

			for _, want := range tt.wantIn {
				if !strings.Contains(prompt, want) {
					t.Errorf("prompt should contain %q", want)
				}
			}

			for _, notWant := range tt.wantOut {
				if strings.Contains(prompt, notWant) {
					t.Errorf("prompt should not contain %q", notWant)
				}
			}
		})
	}
}

func TestBuildPrompt(t *testing.T) {
	prompt := BuildPrompt(ScopeFull, "Title", "Body", "Diff")
	if !strings.Contains(prompt, "Title") {
		t.Error("prompt should contain title")
	}
	if !strings.Contains(prompt, "Body") {
		t.Error("prompt should contain body")
	}
	if !strings.Contains(prompt, "Diff") {
		t.Error("prompt should contain diff")
	}
}

func TestExtractVerdict(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name: "explicit approve",
			content: `## Summary
Good changes.

## Verdict
APPROVE

The code looks good.`,
			want: "APPROVE",
		},
		{
			name: "explicit request changes",
			content: `## Summary
Needs work.

## Verdict
REQUEST_CHANGES

Please fix the security issue.`,
			want: "REQUEST_CHANGES",
		},
		{
			name: "explicit comment",
			content: `## Summary
Some observations.

## Verdict
COMMENT

No major issues.`,
			want: "COMMENT",
		},
		{
			name: "request changes with space",
			content: `## Verdict
REQUEST CHANGES`,
			want: "REQUEST_CHANGES",
		},
		{
			name:    "approve anywhere",
			content: "I APPROVE this PR.",
			want:    "APPROVE",
		},
		{
			name:    "request changes anywhere",
			content: "I REQUEST_CHANGES on this PR.",
			want:    "REQUEST_CHANGES",
		},
		{
			name:    "no verdict defaults to comment",
			content: "This is just a review without a clear verdict.",
			want:    "COMMENT",
		},
		{
			name:    "do not approve",
			content: "I DO NOT APPROVE of this approach.",
			want:    "COMMENT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractVerdict(tt.content)
			if got != tt.want {
				t.Errorf("extractVerdict() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Model != DefaultModel {
		t.Errorf("Model = %s, want %s", cfg.Model, DefaultModel)
	}
	if cfg.Scope != ScopeFull {
		t.Errorf("Scope = %s, want %s", cfg.Scope, ScopeFull)
	}
	if cfg.MaxTokens <= 0 {
		t.Errorf("MaxTokens = %d, want > 0", cfg.MaxTokens)
	}
}

func TestGetProviderConfig_MissingAPIKey(t *testing.T) {
	// This test verifies error handling when API keys are missing
	// We can't easily test successful cases without real API keys

	tests := []struct {
		model   string
		wantErr string
	}{
		{"claude-sonnet-4", "ANTHROPIC_API_KEY"},
		{"gpt-4o", "OPENAI_API_KEY"},
		{"gemini-pro", "GEMINI_API_KEY"},
		{"grok-2", "XAI_API_KEY"},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			_, err := getProviderConfig(tt.model)
			if err == nil {
				t.Skip("API key is set, skipping error test")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error should mention %s, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestGetProviderConfig_Ollama(t *testing.T) {
	// Ollama doesn't require API key
	cfg, err := getProviderConfig("llama3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Provider != "ollama" {
		t.Errorf("Provider = %s, want ollama", cfg.Provider)
	}
}
