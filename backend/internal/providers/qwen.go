package providers

import (
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type QwenProvider struct {
	BaseProvider
}

func NewQwenProvider() *QwenProvider {
	return &QwenProvider{
		BaseProvider: BaseProvider{
			Name:    "qwen",
			BaseURL: "https://dashscope.aliyuncs.com",
		},
	}
}

func (p *QwenProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	// Prefer OAuth token (Qwen Code uses OAuth)
	if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	} else if account.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	}
	return nil
}

// GetValidationEndpoint returns the endpoint for credential validation
func (p *QwenProvider) GetValidationEndpoint() string {
	return "/compatible-mode/v1/models"
}

// GetRateLimitHeaders returns Qwen-specific rate limit header names
func (p *QwenProvider) GetRateLimitHeaders() RateLimitHeaderConfig {
	return RateLimitHeaderConfig{
		RequestsLimit:     "x-ratelimit-limit-requests",
		RequestsRemaining: "x-ratelimit-remaining-requests",
		TokensLimit:       "x-ratelimit-limit-tokens",
		TokensRemaining:   "x-ratelimit-remaining-tokens",
	}
}

// NeedsOAuth indicates if this provider primarily uses OAuth
func (p *QwenProvider) NeedsOAuth() bool {
	return true
}

// GetOAuthConfig returns OAuth configuration for Qwen
func (p *QwenProvider) GetOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		TokenFilePath: "~/.qwen/credentials.json",
	}
}
