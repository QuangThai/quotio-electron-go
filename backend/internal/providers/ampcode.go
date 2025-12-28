package providers

import (
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type AmpcodeProvider struct {
	BaseProvider
}

func NewAmpcodeProvider() *AmpcodeProvider {
	return &AmpcodeProvider{
		BaseProvider: BaseProvider{
			Name:    "ampcode",
			BaseURL: "https://ampcode.com",
		},
	}
}

func (p *AmpcodeProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	}
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
// Ampcode uses OpenAI-compatible endpoints
func (p *AmpcodeProvider) GetValidationEndpoint() string {
	return "/v1/models"
}

// GetRateLimitHeaders returns Ampcode-specific rate limit header names
// Since it's a proxy, it might use standard OpenAI headers
func (p *AmpcodeProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
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
func (p *AmpcodeProvider) NeedsOAuth() bool {
	return false
}

// GetOAuthConfig returns OAuth configuration for Ampcode
func (p *AmpcodeProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.config/ampcode/credentials.json",
	}
}
