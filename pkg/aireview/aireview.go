package aireview

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/plexusone/omnillm"
)

// DefaultModel is the default LLM model to use.
const DefaultModel = "claude-sonnet-4"

// Config holds configuration for AI reviews.
type Config struct {
	// Model is the LLM model to use (e.g., "claude-sonnet-4", "gpt-4o").
	Model string
	// Scope is the review focus (full, security, style, performance).
	Scope Scope
	// MaxTokens is the maximum response length.
	MaxTokens int
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Model:     DefaultModel,
		Scope:     ScopeFull,
		MaxTokens: 4096,
	}
}

// ReviewInput contains the PR information for review.
type ReviewInput struct {
	Title string
	Body  string
	Diff  string
}

// ReviewOutput contains the AI-generated review.
type ReviewOutput struct {
	// Content is the full review text.
	Content string
	// Verdict is the recommended action (APPROVE, COMMENT, REQUEST_CHANGES).
	Verdict string
	// TokensUsed is the total tokens consumed.
	TokensUsed int
}

// Reviewer performs AI-powered code reviews.
type Reviewer struct {
	client *omnillm.ChatClient
	config Config
}

// NewReviewer creates a new AI reviewer with the given configuration.
func NewReviewer(cfg Config) (*Reviewer, error) {
	// Determine provider from model name
	providerCfg, err := getProviderConfig(cfg.Model)
	if err != nil {
		return nil, err
	}

	client, err := omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{providerCfg},
	})
	if err != nil {
		return nil, fmt.Errorf("creating LLM client: %w", err)
	}

	return &Reviewer{
		client: client,
		config: cfg,
	}, nil
}

// Close releases resources held by the reviewer.
func (r *Reviewer) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// Review performs an AI code review on the given PR.
func (r *Reviewer) Review(ctx context.Context, input ReviewInput) (*ReviewOutput, error) {
	// Build the prompt
	userPrompt := BuildUserPrompt(r.config.Scope, input.Title, input.Body, input.Diff)
	systemPrompt := BuildSystemPrompt()

	// Create the request
	maxTokens := r.config.MaxTokens
	temperature := 0.3 // Lower temperature for more focused reviews

	resp, err := r.client.CreateChatCompletion(ctx, &omnillm.ChatCompletionRequest{
		Model: r.config.Model,
		Messages: []omnillm.Message{
			{Role: omnillm.RoleSystem, Content: systemPrompt},
			{Role: omnillm.RoleUser, Content: userPrompt},
		},
		MaxTokens:   &maxTokens,
		Temperature: &temperature,
	})
	if err != nil {
		return nil, fmt.Errorf("LLM request failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from LLM")
	}

	content := resp.Choices[0].Message.Content
	verdict := extractVerdict(content)

	return &ReviewOutput{
		Content:    content,
		Verdict:    verdict,
		TokensUsed: resp.Usage.TotalTokens,
	}, nil
}

// getProviderConfig returns the provider configuration based on model name.
func getProviderConfig(model string) (omnillm.ProviderConfig, error) {
	model = strings.ToLower(model)

	switch {
	case strings.HasPrefix(model, "claude") || strings.HasPrefix(model, "anthropic"):
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return omnillm.ProviderConfig{}, fmt.Errorf("ANTHROPIC_API_KEY environment variable is required for Claude models")
		}
		return omnillm.ProviderConfig{
			Provider: omnillm.ProviderNameAnthropic,
			APIKey:   apiKey,
		}, nil

	case strings.HasPrefix(model, "gpt") || strings.HasPrefix(model, "openai"):
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return omnillm.ProviderConfig{}, fmt.Errorf("OPENAI_API_KEY environment variable is required for OpenAI models")
		}
		return omnillm.ProviderConfig{
			Provider: omnillm.ProviderNameOpenAI,
			APIKey:   apiKey,
		}, nil

	case strings.HasPrefix(model, "gemini") || strings.HasPrefix(model, "google"):
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return omnillm.ProviderConfig{}, fmt.Errorf("GEMINI_API_KEY environment variable is required for Gemini models")
		}
		return omnillm.ProviderConfig{
			Provider: omnillm.ProviderNameGemini,
			APIKey:   apiKey,
		}, nil

	case strings.HasPrefix(model, "grok") || strings.HasPrefix(model, "xai"):
		apiKey := os.Getenv("XAI_API_KEY")
		if apiKey == "" {
			return omnillm.ProviderConfig{}, fmt.Errorf("XAI_API_KEY environment variable is required for Grok models")
		}
		return omnillm.ProviderConfig{
			Provider: omnillm.ProviderNameXAI,
			APIKey:   apiKey,
		}, nil

	case strings.HasPrefix(model, "llama") || strings.HasPrefix(model, "mistral") || strings.HasPrefix(model, "ollama"):
		baseURL := os.Getenv("OLLAMA_BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		return omnillm.ProviderConfig{
			Provider: omnillm.ProviderNameOllama,
			BaseURL:  baseURL,
		}, nil

	default:
		// Default to Anthropic for unknown models
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return omnillm.ProviderConfig{}, fmt.Errorf("ANTHROPIC_API_KEY environment variable is required (defaulting to Anthropic for model %q)", model)
		}
		return omnillm.ProviderConfig{
			Provider: omnillm.ProviderNameAnthropic,
			APIKey:   apiKey,
		}, nil
	}
}

// extractVerdict extracts the verdict from the review content.
func extractVerdict(content string) string {
	content = strings.ToUpper(content)

	// Look for explicit verdict markers
	if strings.Contains(content, "## VERDICT") {
		// Find the verdict section and extract the action
		idx := strings.Index(content, "## VERDICT")
		section := content[idx:]

		if strings.Contains(section, "APPROVE") {
			return "APPROVE"
		}
		if strings.Contains(section, "REQUEST_CHANGES") || strings.Contains(section, "REQUEST CHANGES") {
			return "REQUEST_CHANGES"
		}
		if strings.Contains(section, "COMMENT") {
			return "COMMENT"
		}
	}

	// Fallback: look for keywords anywhere
	if strings.Contains(content, "REQUEST_CHANGES") || strings.Contains(content, "REQUEST CHANGES") {
		return "REQUEST_CHANGES"
	}
	if strings.Contains(content, "APPROVE") && !strings.Contains(content, "DO NOT APPROVE") {
		return "APPROVE"
	}

	// Default to COMMENT if verdict is unclear
	return "COMMENT"
}

// GetModelFromEnv returns the model from environment variable or default.
func GetModelFromEnv() string {
	if model := os.Getenv("OMNILLM_MODEL"); model != "" {
		return model
	}
	return DefaultModel
}

// GetScopeFromEnv returns the scope from environment variable or default.
func GetScopeFromEnv() Scope {
	if scope := os.Getenv("ACR_REVIEW_SCOPE"); scope != "" && IsValidScope(scope) {
		return Scope(scope)
	}
	return ScopeFull
}
