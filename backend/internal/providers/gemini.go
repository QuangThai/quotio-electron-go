package providers

import (
	"context"
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type GeminiProvider struct {
	BaseProvider
}

func NewGeminiProvider() *GeminiProvider {
	return &GeminiProvider{
		BaseProvider: BaseProvider{
			Name:    "gemini",
			BaseURL: "https://generativelanguage.googleapis.com",
		},
	}
}

func (p *GeminiProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	// Prefer OAuth token (Gemini CLI uses OAuth)
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		// Gemini uses API key as query parameter
		q := req.URL.Query()
		q.Set("key", account.APIKey)
		req.URL.RawQuery = q.Encode()
	}
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
func (p *GeminiProvider) GetValidationEndpoint() string {
	return "/v1beta/models"
}

// GetRateLimitHeaders returns Gemini-specific rate limit header names
func (p *GeminiProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
	return RateLimitHeaderConfig{
		RequestsLimit:     "x-ratelimit-limit",
		RequestsRemaining: "x-ratelimit-remaining",
		RequestsReset:     "x-ratelimit-reset",
	}
}

// NeedsOAuth indicates if this provider primarily uses OAuth
func (p *GeminiProvider) NeedsOAuth() bool {
	return true // Gemini CLI uses OAuth
}

// GetOAuthConfig returns OAuth configuration for Gemini
func (p *GeminiProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.gemini/credentials.json",
	}
}

func (p *GeminiProvider) FetchQuota(ctx context.Context, account *storage.Account) (int64, int64, error) {
	// For Gemini, we don't have a dedicated 'get quota' endpoint.
	// We rely on the last seen headers or we could make a dummy request.
	// For now, we'll return 0, 0 to signal 'unknown' or 'use cached' 
	// unless we implement a specific management API call.
	
	// If we have an OAuth token, we might be able to call a broader Google Cloud API,
	// but for standard Gemini API key, there's no direct quota endpoint.
	return 0, 0, nil
}
