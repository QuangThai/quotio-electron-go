package providers

import (
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type CopilotProvider struct {
	BaseProvider
}

func NewCopilotProvider() *CopilotProvider {
	return &CopilotProvider{
		BaseProvider: BaseProvider{
			Name:    "copilot",
			BaseURL: "https://api.github.com",
		},
	}
}

func (p *CopilotProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	// GitHub Copilot uses OAuth tokens from GitHub
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
func (p *CopilotProvider) GetValidationEndpoint() string {
	return "/user"
}

// GetRateLimitHeaders returns GitHub-specific rate limit header names
func (p *CopilotProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
	return RateLimitHeaderConfig{
		RequestsLimit:     "x-ratelimit-limit",
		RequestsRemaining: "x-ratelimit-remaining",
		RequestsReset:     "x-ratelimit-reset",
	}
}

// NeedsOAuth indicates if this provider primarily uses OAuth
func (p *CopilotProvider) NeedsOAuth() bool {
	return true
}

// GetOAuthConfig returns OAuth configuration for GitHub Copilot
func (p *CopilotProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.config/github-copilot/hosts.json",
	}
}
