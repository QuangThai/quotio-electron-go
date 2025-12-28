package providers

import (
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type ClaudeProvider struct {
	BaseProvider
}

func NewClaudeProvider() *ClaudeProvider {
	return &ClaudeProvider{
		BaseProvider: BaseProvider{
			Name:    "claude",
			BaseURL: "https://api.anthropic.com",
		},
	}
}

func (p *ClaudeProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	// Prefer OAuth token over API key (Claude Code uses OAuth)
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		req.Header.Set("x-api-key", account.APIKey)
	}
	req.Header.Set("anthropic-version", "2023-06-01")
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
func (p *ClaudeProvider) GetValidationEndpoint() string {
	return "/v1/messages"
}

// GetRateLimitHeaders returns Claude/Anthropic-specific rate limit header names
func (p *ClaudeProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
	return RateLimitHeaderConfig{
		RequestsLimit:        "anthropic-ratelimit-requests-limit",
		RequestsRemaining:    "anthropic-ratelimit-requests-remaining",
		RequestsReset:        "anthropic-ratelimit-requests-reset",
		TokensLimit:          "anthropic-ratelimit-tokens-limit",
		TokensRemaining:      "anthropic-ratelimit-tokens-remaining",
		TokensReset:          "anthropic-ratelimit-tokens-reset",
		InputTokensLimit:     "anthropic-ratelimit-input-tokens-limit",
		InputTokensRemaining: "anthropic-ratelimit-input-tokens-remaining",
		OutputTokensLimit:    "anthropic-ratelimit-output-tokens-limit",
		OutputTokensRemaining: "anthropic-ratelimit-output-tokens-remaining",
	}
}

// NeedsOAuth indicates if this provider primarily uses OAuth
func (p *ClaudeProvider) NeedsOAuth() bool {
	return true // Claude Code uses OAuth
}

// GetOAuthConfig returns OAuth configuration for Claude
func (p *ClaudeProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.claude/credentials.json",
	}
}
