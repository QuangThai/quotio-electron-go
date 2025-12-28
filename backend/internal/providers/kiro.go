package providers

import (
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type KiroProvider struct {
	BaseProvider
}

func NewKiroProvider() *KiroProvider {
	return &KiroProvider{
		BaseProvider: BaseProvider{
			Name:    "kiro",
			BaseURL: "https://api.kiro.ai",
		},
	}
}

func (p *KiroProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	// Kiro (Amazon CodeWhisperer) uses OAuth tokens
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	}
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
func (p *KiroProvider) GetValidationEndpoint() string {
	return "/v1/models"
}

// GetRateLimitHeaders returns Kiro-specific rate limit header names
func (p *KiroProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
	return RateLimitHeaderConfig{
		RequestsLimit:     "x-ratelimit-limit",
		RequestsRemaining: "x-ratelimit-remaining",
	}
}

// NeedsOAuth indicates if this provider primarily uses OAuth
func (p *KiroProvider) NeedsOAuth() bool {
	return true
}

// GetOAuthConfig returns OAuth configuration for Kiro
func (p *KiroProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.kiro/credentials.json",
	}
}
