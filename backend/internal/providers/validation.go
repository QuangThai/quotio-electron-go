package providers

import (
	"fmt"
	"net/http"
	"time"
	"quotio-electron-go/backend/internal/storage"
)

// CredentialValidationResult contains validation outcome
type CredentialValidationResult struct {
	IsValid  bool
	Reason   string // "success", "invalid_credentials", "network_error", etc.
	ErrorMsg string
}

// ValidateAccountCredentials tests if an account has working provider credentials
// It performs a minimal, non-destructive check (e.g., list models endpoint)
func ValidateAccountCredentials(account *storage.Account) CredentialValidationResult {
	if account == nil || account.Provider == "" {
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   "empty_account",
			ErrorMsg: "Account is empty",
		}
	}

	provider := GetProviderForAccount(account)
	if provider == nil {
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   "unknown_provider",
			ErrorMsg: fmt.Sprintf("Provider %s not supported", account.Provider),
		}
	}

	// Create minimal test request
	client := &http.Client{Timeout: 5 * time.Second}

	// Build test URL based on provider - use provider's specific validation endpoint
	testURL := fmt.Sprintf("%s%s", provider.GetBaseURL(), provider.GetValidationEndpoint())
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   "request_creation_error",
			ErrorMsg: err.Error(),
		}
	}

	// Apply provider auth
	if err := provider.AuthenticateRequest(req, account); err != nil {
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   "auth_header_error",
			ErrorMsg: err.Error(),
		}
	}

	// Make test request
	resp, err := client.Do(req)
	if err != nil {
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   "network_error",
			ErrorMsg: err.Error(),
		}
	}
	defer resp.Body.Close()

	// Interpret response
	switch {
	case resp.StatusCode == 200:
		return CredentialValidationResult{IsValid: true, Reason: "success"}
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   "invalid_credentials",
			ErrorMsg: "Provider rejected credentials (401/403)",
		}
	case resp.StatusCode == 429:
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   "rate_limited",
			ErrorMsg: "Provider is rate limiting",
		}
	default:
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   fmt.Sprintf("http_%d", resp.StatusCode),
			ErrorMsg: fmt.Sprintf("Provider returned %d", resp.StatusCode),
		}
	}
}
