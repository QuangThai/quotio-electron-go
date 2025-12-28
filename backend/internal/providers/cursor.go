package providers

import (
	"context"
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type CursorProvider struct {
	BaseProvider
}

func NewCursorProvider() *CursorProvider {
	return &CursorProvider{
		BaseProvider: BaseProvider{
			Name:    "cursor",
			BaseURL: "https://api.cursor.sh",
		},
	}
}

func (p *CursorProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	if account.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	}
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
func (p *CursorProvider) GetValidationEndpoint() string {
	return "/v1/models"
}

// GetRateLimitHeaders returns Cursor-specific rate limit header names
func (p *CursorProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
	return RateLimitHeaderConfig{
		RequestsLimit:     "x-ratelimit-limit-requests",
		RequestsRemaining: "x-ratelimit-remaining-requests",
		RequestsReset:     "x-ratelimit-reset-requests",
	}
}

// NeedsOAuth indicates if this provider primarily uses OAuth
func (p *CursorProvider) NeedsOAuth() bool {
	return false
}

// GetOAuthConfig returns OAuth configuration for Cursor
func (p *CursorProvider) GetOAuthConfig() *OAuthConfig {
	return nil
}

func (p *CursorProvider) FetchQuota(ctx context.Context, account *storage.Account) (int64, int64, error) {
	// Cursor doesn't provide a public quota endpoint that we can easily use here.
	return 0, 0, nil
}
