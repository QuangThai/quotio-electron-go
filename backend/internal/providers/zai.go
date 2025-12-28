package providers

import (
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type ZAIProvider struct {
	BaseProvider
}

func NewZAIProvider() *ZAIProvider {
	return &ZAIProvider{
		BaseProvider: BaseProvider{
			Name:    "z.ai",
			BaseURL: "https://api.z.ai/api/paas/v4",
		},
	}
}

func (p *ZAIProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	}
	return nil
}

// GetValidationEndpoint for Z.AI is just /models since /v4 is in base URL
func (p *ZAIProvider) GetValidationEndpoint() string {
	return "/models"
}

// GetRateLimitHeaders returns Z.AI-specific rate limit header names
func (p *ZAIProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
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
func (p *ZAIProvider) NeedsOAuth() bool {
	return false
}

// GetOAuthConfig returns OAuth configuration for Z.AI
func (p *ZAIProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.config/zai/credentials.json",
	}
}
