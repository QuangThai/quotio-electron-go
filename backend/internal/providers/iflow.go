package providers

import (
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type IFlowProvider struct {
	BaseProvider
}

func NewIFlowProvider() *IFlowProvider {
	return &IFlowProvider{
		BaseProvider: BaseProvider{
			Name:    "iflow",
			BaseURL: "https://api.iflow.ai",
		},
	}
}

func (p *IFlowProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	}
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
func (p *IFlowProvider) GetValidationEndpoint() string {
	return "/v1/models"
}

// GetRateLimitHeaders returns iFlow-specific rate limit header names
func (p *IFlowProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
	return RateLimitHeaderConfig{
		RequestsLimit:     "x-ratelimit-limit-requests",
		RequestsRemaining: "x-ratelimit-remaining-requests",
	}
}

// NeedsOAuth indicates if this provider primarily uses OAuth
func (p *IFlowProvider) NeedsOAuth() bool {
	return true
}

// GetOAuthConfig returns OAuth configuration for iFlow
func (p *IFlowProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.iflow/credentials.json",
	}
}
