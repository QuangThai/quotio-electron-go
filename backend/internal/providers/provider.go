package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"quotio-electron-go/backend/internal/storage"
)

type Provider interface {
	GetName() string
	GetBaseURL() string
	AuthenticateRequest(req *http.Request, account *storage.Account) error
	ParseQuotaFromResponse(resp *http.Response) (int64, error) // Parse from headers only
	ParseQuotaFromBody(body []byte) (int64, error)              // Parse quota from buffered body
	DetectRateLimit(resp *http.Response) bool
	GetValidationEndpoint() string
	FetchQuota(ctx context.Context, account *storage.Account) (int64, int64, error) // Returns (used, limit, error)
}

type BaseProvider struct {
	Name    string
	BaseURL string
}

func (p *BaseProvider) GetName() string {
	return p.Name
}

func (p *BaseProvider) GetBaseURL() string {
	return p.BaseURL
}

func (p *BaseProvider) AuthenticateRequest(req *http.Request, account *storage.Account) error {
	if account.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	} else if account.OAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+account.OAuthToken)
	}
	return nil
}

func (p *BaseProvider) ParseQuotaFromResponse(resp *http.Response) (int64, error) {
	// Default implementation - parse from response headers only
	// Do NOT consume the response body here as it breaks streaming responses (SSE)
	// Instead, providers should override this method to parse headers if needed
	// For body parsing, use the new ParseQuotaFromBody helper in proxy/server.go
	return 0, nil
}

func (p *BaseProvider) ParseQuotaFromBody(body []byte) (int64, error) {
	// Default implementation - parse from buffered response body
	// This is called with a buffered copy of the body, safe for non-streaming responses
	// Streaming responses (SSE) should override to return 0 or parse from headers only
	
	// Try to parse JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	// Look for usage fields (provider-specific)
	if usage, ok := data["usage"].(map[string]interface{}); ok {
		if tokens, ok := usage["total_tokens"].(float64); ok {
			return int64(tokens), nil
		}
	}

	return 0, nil
}

func (p *BaseProvider) DetectRateLimit(resp *http.Response) bool {
	return resp.StatusCode == 429 || resp.StatusCode == 403
}

func (p *BaseProvider) GetValidationEndpoint() string {
	return "/v1/models"
}

func (p *BaseProvider) FetchQuota(ctx context.Context, account *storage.Account) (int64, int64, error) {
	// Default implementation returns 0, 0
	return 0, 0, nil
}

