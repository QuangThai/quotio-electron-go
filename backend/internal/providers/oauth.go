package providers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// OAuthCredentials represents stored OAuth credentials
type OAuthCredentials struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
}

// OAuthConfig holds provider-specific OAuth configuration
type OAuthConfig struct {
	AuthURL       string
	TokenURL      string
	ClientID      string
	Scopes        []string
	TokenFilePath string // Path to CLI auth file
}

// RateLimitHeaderConfig defines provider-specific rate limit header names
type RateLimitHeaderConfig struct {
	RequestsLimit     string
	RequestsRemaining string
	RequestsReset     string
	TokensLimit       string
	TokensRemaining   string
	TokensReset       string
	InputTokensLimit  string
	InputTokensRemaining string
	OutputTokensLimit    string
	OutputTokensRemaining string
}

// LoadOAuthFromFile loads OAuth credentials from CLI auth files
func LoadOAuthFromFile(path string) (*OAuthCredentials, error) {
	expandedPath := expandPath(path)
	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, err
	}

	var creds OAuthCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		// Try alternate JSON structures
		var altCreds map[string]interface{}
		if err := json.Unmarshal(data, &altCreds); err != nil {
			return nil, err
		}
		
		// Parse common OAuth response formats
		if token, ok := altCreds["access_token"].(string); ok {
			creds.AccessToken = token
		} else if token, ok := altCreds["token"].(string); ok {
			creds.AccessToken = token
		}
		
		if refresh, ok := altCreds["refresh_token"].(string); ok {
			creds.RefreshToken = refresh
		}
		
		if tokenType, ok := altCreds["token_type"].(string); ok {
			creds.TokenType = tokenType
		}
	}

	return &creds, nil
}

// DetectProviderCredentials detects OAuth credentials from common paths
func DetectProviderCredentials(provider string) (*OAuthCredentials, error) {
	paths := getProviderAuthPaths(provider)
	for _, path := range paths {
		if creds, err := LoadOAuthFromFile(path); err == nil && creds.AccessToken != "" {
			return creds, nil
		}
	}
	return nil, nil
}

// getProviderAuthPaths returns common credential file paths for each provider
func getProviderAuthPaths(provider string) []string {
	home, _ := os.UserHomeDir()
	switch provider {
	case "claude":
		return []string{
			filepath.Join(home, ".claude", "credentials.json"),
			filepath.Join(home, ".config", "claude", "credentials.json"),
			filepath.Join(home, ".claude.json"),
		}
	case "openai", "codex":
		return []string{
			filepath.Join(home, ".openai", "credentials.json"),
			filepath.Join(home, ".config", "openai", "credentials.json"),
			filepath.Join(home, ".codex", "credentials.json"),
		}
	case "gemini":
		return []string{
			filepath.Join(home, ".gemini", "credentials.json"),
			filepath.Join(home, ".config", "gemini", "credentials"),
			filepath.Join(home, ".gemini_cli", "credentials.json"),
		}
	case "qwen":
		return []string{
			filepath.Join(home, ".qwen", "credentials.json"),
			filepath.Join(home, ".config", "qwen", "credentials.json"),
		}
	case "iflow":
		return []string{
			filepath.Join(home, ".iflow", "credentials.json"),
		}
	case "antigravity":
		return []string{
			filepath.Join(home, ".antigravity", "credentials.json"),
			filepath.Join(home, ".config", "antigravity", "credentials.json"),
		}
	case "copilot":
		return []string{
			filepath.Join(home, ".config", "github-copilot", "hosts.json"),
			filepath.Join(home, ".github-copilot", "hosts.json"),
		}
	case "kiro":
		return []string{
			filepath.Join(home, ".kiro", "credentials.json"),
			filepath.Join(home, ".aws", "credentials"),
		}
	case "ampcode":
		return []string{
			filepath.Join(home, ".ampcode", "credentials.json"),
			filepath.Join(home, ".config", "ampcode", "credentials.json"),
		}
	case "z.ai":
		return []string{
			filepath.Join(home, ".zai", "credentials.json"),
			filepath.Join(home, ".config", "zai", "credentials.json"),
		}
	}
	return nil
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

// GetAllProviderCredentialPaths returns all paths to check for credentials
func GetAllProviderCredentialPaths() map[string][]string {
	providers := []string{"claude", "openai", "gemini", "qwen", "iflow", "antigravity", "copilot", "kiro", "ampcode", "z.ai"}
	result := make(map[string][]string)
	for _, p := range providers {
		result[p] = getProviderAuthPaths(p)
	}
	return result
}

// DetectAllCredentials scans for all available provider credentials
func DetectAllCredentials() map[string]*OAuthCredentials {
	providers := []string{"claude", "openai", "gemini", "qwen", "iflow", "antigravity", "copilot", "kiro", "ampcode", "z.ai"}
	result := make(map[string]*OAuthCredentials)
	
	for _, provider := range providers {
		if creds, err := DetectProviderCredentials(provider); err == nil && creds != nil {
			result[provider] = creds
		}
	}
	
	return result
}
