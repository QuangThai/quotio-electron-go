package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"quotio-electron-go/backend/internal/storage"
	"time"
)

// ProviderClient wraps provider-specific client logic
type ProviderClient struct {
	account *storage.Account
	client  *http.Client
}

// NewProviderClient creates a client for a given account
func NewProviderClient(account *storage.Account) (*ProviderClient, error) {
	if account == nil || account.Provider == "" {
		return nil, errors.New("empty account or provider")
	}

	return &ProviderClient{
		account: account,
		client:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// ValidateCredentials tests if credentials work
func (pc *ProviderClient) ValidateCredentials(ctx context.Context) (bool, string, error) {
	if pc.account == nil || pc.account.Provider == "" {
		return false, "empty_account", errors.New("account is empty")
	}

	provider := GetProviderForAccount(pc.account)
	if provider == nil {
		return false, "unknown_provider", fmt.Errorf("provider %s not supported", pc.account.Provider)
	}

	// Create test request to validation endpoint
	testURL := fmt.Sprintf("%s%s", provider.GetBaseURL(), provider.GetValidationEndpoint())
	req, err := http.NewRequestWithContext(ctx, "GET", testURL, nil)
	if err != nil {
		return false, "request_creation_error", err
	}

	// Apply authentication
	if err := provider.AuthenticateRequest(req, pc.account); err != nil {
		return false, "auth_header_error", err
	}

	// Make request
	resp, err := pc.client.Do(req)
	if err != nil {
		return false, "network_error", err
	}
	defer resp.Body.Close()

	// Parse response
	switch {
	case resp.StatusCode == 200:
		return true, "success", nil
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		return false, "invalid_credentials", fmt.Errorf("provider returned %d", resp.StatusCode)
	case resp.StatusCode == 429:
		return false, "rate_limited", fmt.Errorf("provider rate limited")
	default:
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Sprintf("http_%d", resp.StatusCode), fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
}

// GetRateLimits fetches current rate limits from provider response headers
func (pc *ProviderClient) GetRateLimits(resp *http.Response) (*RateLimitInfo, error) {
	if resp == nil {
		return nil, errors.New("nil response")
	}

	provider := GetProviderForAccount(pc.account)
	if provider == nil {
		return nil, fmt.Errorf("provider %s not supported", pc.account.Provider)
	}

	// Parse rate limit headers based on provider
	return parseProviderRateLimits(pc.account.Provider, resp.Header), nil
}

// parseProviderRateLimits parses rate limit headers from different providers
func parseProviderRateLimits(provider string, headers http.Header) *RateLimitInfo {
	info := &RateLimitInfo{}

	switch provider {
	case "claude":
		// Anthropic headers
		if val := headers.Get("anthropic-ratelimit-requests-limit"); val != "" {
			info.RequestsLimit = parseInt64(val)
		}
		if val := headers.Get("anthropic-ratelimit-requests-remaining"); val != "" {
			info.RequestsRemaining = parseInt64(val)
		}
		if val := headers.Get("anthropic-ratelimit-requests-reset"); val != "" {
			info.RequestsReset = parseTime(val)
		}
		if val := headers.Get("anthropic-ratelimit-tokens-limit"); val != "" {
			info.TokensLimit = parseInt64(val)
		}
		if val := headers.Get("anthropic-ratelimit-tokens-remaining"); val != "" {
			info.TokensRemaining = parseInt64(val)
		}
		if val := headers.Get("anthropic-ratelimit-tokens-reset"); val != "" {
			info.TokensReset = parseTime(val)
		}
		// Input/output tokens (Claude specific)
		if val := headers.Get("anthropic-ratelimit-input-tokens-limit"); val != "" {
			info.InputTokensLimit = parseInt64(val)
		}
		if val := headers.Get("anthropic-ratelimit-input-tokens-remaining"); val != "" {
			info.InputTokensRemaining = parseInt64(val)
		}
		if val := headers.Get("anthropic-ratelimit-input-tokens-reset"); val != "" {
			info.InputTokensReset = parseTime(val)
		}
		if val := headers.Get("anthropic-ratelimit-output-tokens-limit"); val != "" {
			info.OutputTokensLimit = parseInt64(val)
		}
		if val := headers.Get("anthropic-ratelimit-output-tokens-remaining"); val != "" {
			info.OutputTokensRemaining = parseInt64(val)
		}
		if val := headers.Get("anthropic-ratelimit-output-tokens-reset"); val != "" {
			info.OutputTokensReset = parseTime(val)
		}

	case "openai":
		// OpenAI headers
		if val := headers.Get("x-ratelimit-limit-requests"); val != "" {
			info.RequestsLimit = parseInt64(val)
		}
		if val := headers.Get("x-ratelimit-remaining-requests"); val != "" {
			info.RequestsRemaining = parseInt64(val)
		}
		if val := headers.Get("x-ratelimit-reset-requests"); val != "" {
			info.RequestsReset = parseTime(val)
		}
		if val := headers.Get("x-ratelimit-limit-tokens"); val != "" {
			info.TokensLimit = parseInt64(val)
		}
		if val := headers.Get("x-ratelimit-remaining-tokens"); val != "" {
			info.TokensRemaining = parseInt64(val)
		}
		if val := headers.Get("x-ratelimit-reset-tokens"); val != "" {
			info.TokensReset = parseTime(val)
		}

	case "gemini", "vertex":
		// Google/Vertex headers (if available)
		if val := headers.Get("X-Server-Endpoint-Quota"); val != "" {
			// Parse quota from response body instead for Gemini
			// We'll handle this in the response body parsing
		}
		// Gemini does not provide rate limit headers
		info.RequestsLimit = -1
		info.TokensLimit = -1

	case "qwen":
		// Qwen headers (similar to OpenAI)
		if val := headers.Get("x-ratelimit-limit-requests"); val != "" {
			info.RequestsLimit = parseInt64(val)
		}
		if val := headers.Get("x-ratelimit-remaining-requests"); val != "" {
			info.RequestsRemaining = parseInt64(val)
		}

	case "antigravity", "iflow", "kiro", "copilot":
		// These providers may not expose rate limit headers
		// We'll rely on response body parsing
	}

	return info
}

// ParseRateLimitsFromResponse parses rate limits from full response
func (pc *ProviderClient) ParseRateLimitsFromResponse(resp *http.Response) (*RateLimitInfo, error) {
	if resp == nil {
		return nil, errors.New("nil response")
	}

	// Parse headers
	info := parseProviderRateLimits(pc.account.Provider, resp.Header)

	// If no rate limit headers found, try parsing from response body
	if info.TokensLimit == 0 && info.RequestsLimit == 0 {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			info = parseRateLimitsFromBody(pc.account.Provider, body)
		}
	}

	return info, nil
}

// parseRateLimitsFromBody attempts to extract rate limits from response JSON
func parseRateLimitsFromBody(provider string, body []byte) *RateLimitInfo {
	info := &RateLimitInfo{}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return info
	}

	// Try common rate limit fields
	if usage, ok := data["usage"].(map[string]interface{}); ok {
		if tokens, ok := usage["total_tokens"].(float64); ok {
			info.TokensUsed = int64(tokens)
		}
	}

	// Provider-specific body parsing
	switch provider {
	case "openai":
		if usage, ok := data["usage"].(map[string]interface{}); ok {
			if total, ok := usage["total_tokens"].(float64); ok {
				info.TokensUsed = int64(total)
			}
		}

	case "claude":
		if usage, ok := data["usage"].(map[string]interface{}); ok {
			if input, ok := usage["input_tokens"].(float64); ok {
				info.InputTokensUsed = int64(input)
			}
			if output, ok := usage["output_tokens"].(float64); ok {
				info.OutputTokensUsed = int64(output)
			}
			info.TokensUsed = info.InputTokensUsed + info.OutputTokensUsed
		}

	case "gemini":
		if usage, ok := data["usageMetadata"].(map[string]interface{}); ok {
			if total, ok := usage["totalTokenCount"].(float64); ok {
				info.TokensUsed = int64(total)
			}
			if prompt, ok := usage["promptTokenCount"].(float64); ok {
				info.InputTokensUsed = int64(prompt)
			}
			if candidates, ok := usage["candidatesTokenCount"].(float64); ok {
				info.OutputTokensUsed = int64(candidates)
			}
		}
	}

	return info
}

// RateLimitInfo contains rate limit information from provider
type RateLimitInfo struct {
	RequestsLimit        int64     `json:"requests_limit"`
	RequestsRemaining     int64     `json:"requests_remaining"`
	RequestsReset        time.Time `json:"requests_reset"`
	TokensLimit          int64     `json:"tokens_limit"`
	TokensRemaining       int64     `json:"tokens_remaining"`
	TokensReset          time.Time `json:"tokens_reset"`
	InputTokensLimit     int64     `json:"input_tokens_limit"`
	InputTokensRemaining  int64     `json:"input_tokens_remaining"`
	InputTokensReset     time.Time `json:"input_tokens_reset"`
	OutputTokensLimit    int64     `json:"output_tokens_limit"`
	OutputTokensRemaining int64     `json:"output_tokens_remaining"`
	OutputTokensReset   time.Time `json:"output_tokens_reset"`
	TokensUsed           int64     `json:"tokens_used"`
	InputTokensUsed      int64     `json:"input_tokens_used"`
	OutputTokensUsed     int64     `json:"output_tokens_used"`
}

// Helper functions
func parseInt64(s string) int64 {
	var result int64
	fmt.Sscanf(s, "%d", &result)
	return result
}

func parseTime(s string) time.Time {
	// Try RFC3339 format first
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	// Fallback to Unix timestamp
	var sec int64
	fmt.Sscanf(s, "%d", &sec)
	return time.Unix(sec, 0)
}

