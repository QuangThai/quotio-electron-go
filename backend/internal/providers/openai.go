package providers

import (
	"context"
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type OpenAIProvider struct {
	BaseProvider
}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		BaseProvider: BaseProvider{
			Name:    "openai",
			BaseURL: "https://api.openai.com",
		},
	}
}

func (p *OpenAIProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	// Prefer OAuth token over API key (Codex uses OAuth)
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	}
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
func (p *OpenAIProvider) GetValidationEndpoint() string {
	return "/v1/models"
}

// GetRateLimitHeaders returns OpenAI-specific rate limit header names
func (p *OpenAIProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
	return RateLimitHeaderConfig{
		RequestsLimit:     "x-ratelimit-limit-requests",
		RequestsRemaining: "x-ratelimit-remaining-requests",
		RequestsReset:     "x-ratelimit-reset-requests",
		TokensLimit:       "x-ratelimit-limit-tokens",
		TokensRemaining:   "x-ratelimit-remaining-tokens",
		TokensReset:       "x-ratelimit-reset-tokens",
	}
}

// NeedsOAuth indicates if this provider primarily uses OAuth
func (p *OpenAIProvider) NeedsOAuth() bool {
	return true // OpenAI Codex uses OAuth
}

// GetOAuthConfig returns OAuth configuration for OpenAI
func (p *OpenAIProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.openai/credentials.json",
	}
}

func (p *OpenAIProvider) FetchQuota(ctx context.Context, account *storage.Account) (int64, int64, error) {
	// OpenAI standard API doesn't provide a quota endpoint.
	// Usage can be fetched via https://api.openai.com/v1/dashboard/billing/usage
	// but that requires specific permissions/tokens.
	return 0, 0, nil
}
