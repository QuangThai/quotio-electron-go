package providers

import (
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type VertexProvider struct {
	BaseProvider
}

func NewVertexProvider() *VertexProvider {
	return &VertexProvider{
		BaseProvider: BaseProvider{
			Name:    "vertex",
			BaseURL: "https://aiplatform.googleapis.com",
		},
	}
}

func (p *VertexProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	// Vertex AI uses OAuth tokens (Google Cloud)
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	}
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
func (p *VertexProvider) GetValidationEndpoint() string {
	return "/v1/models"
}

// GetRateLimitHeaders returns Vertex AI-specific rate limit header names
func (p *VertexProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
	return RateLimitHeaderConfig{
		RequestsLimit:     "x-ratelimit-limit",
		RequestsRemaining: "x-ratelimit-remaining",
	}
}

// NeedsOAuth indicates if this provider primarily uses OAuth
func (p *VertexProvider) NeedsOAuth() bool {
	return true
}

// GetOAuthConfig returns OAuth configuration for Vertex AI
func (p *VertexProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.config/gcloud/application_default_credentials.json",
	}
}
