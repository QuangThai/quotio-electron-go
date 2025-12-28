package providers

import (
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type AntigravityProvider struct {
	BaseProvider
}

func NewAntigravityProvider() *AntigravityProvider {
	return &AntigravityProvider{
		BaseProvider: BaseProvider{
			Name: "antigravity",
			// Antigravity uses Gemini infrastructure to access Claude models
			// It's Google's Unified Gateway API
			BaseURL: "https://generativelanguage.googleapis.com",
		},
	}
}

func (p *AntigravityProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	// Prefer OAuth token (for cloud proxy)
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		// Use API key as query parameter (same as Gemini)
		q := req.URL.Query()
		q.Set("key", account.APIKey)
		req.URL.RawQuery = q.Encode()
	}
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
// Antigravity validates via Gemini's models endpoint
func (p *AntigravityProvider) GetValidationEndpoint() string {
	return "/v1beta/models"
}

// GetRateLimitHeaders returns Antigravity-specific rate limit header names
func (p *AntigravityProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
	return RateLimitHeaderConfig{
		RequestsLimit:     "x-ratelimit-limit-requests",
		RequestsRemaining: "x-ratelimit-remaining-requests",
	}
}

// NeedsOAuth indicates if this provider primarily uses OAuth
func (p *AntigravityProvider) NeedsOAuth() bool {
	return true
}

// GetOAuthConfig returns OAuth configuration for Antigravity
func (p *AntigravityProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.antigravity/credentials.json",
	}
}
