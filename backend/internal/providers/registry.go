package providers

import (
	"context"
	"fmt"
	"quotio-electron-go/backend/internal/storage"
	"sync"
	"time"
)

var (
	registry     map[string]Provider
	registryOnce sync.Once
)

func init() {
	registryOnce.Do(func() {
		registry = make(map[string]Provider)

		// Register all providers
		RegisterProvider(NewClaudeProvider())
		RegisterProvider(NewOpenAIProvider())
		RegisterProvider(NewGeminiProvider())
		RegisterProvider(NewAntigravityProvider())
		RegisterProvider(NewCopilotProvider())
		RegisterProvider(NewQwenProvider())
		RegisterProvider(NewVertexProvider())
		RegisterProvider(NewIFlowProvider())
		RegisterProvider(NewKiroProvider())
		RegisterProvider(NewAmpcodeProvider())
		RegisterProvider(NewZAIProvider())
		RegisterProvider(NewCursorProvider())
	})
}

func RegisterProvider(provider Provider) {
	registry[provider.GetName()] = provider
}

func GetProvider(name string) Provider {
	return registry[name]
}

func GetProviderForAccount(account *storage.Account) Provider {
	return GetProvider(account.Provider)
}

// CreateProviderClient creates a provider client for a given account
func CreateProviderClient(account *storage.Account) (*ProviderClient, error) {
	return NewProviderClient(account)
}

// ValidateAccountWithClient validates account credentials using provider client
func ValidateAccountWithClient(account *storage.Account) CredentialValidationResult {
	client, err := CreateProviderClient(account)
	if err != nil {
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   "client_creation_failed",
			ErrorMsg: err.Error(),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	isValid, reason, err := client.ValidateCredentials(ctx)
	if err != nil {
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   "validation_error",
			ErrorMsg: err.Error(),
		}
	}

	if !isValid {
		return CredentialValidationResult{
			IsValid:  false,
			Reason:   reason,
			ErrorMsg: "Provider rejected credentials",
		}
	}

	return CredentialValidationResult{IsValid: true, Reason: "success"}
}

// FetchProviderQuota fetches actual quota from provider
func FetchProviderQuota(account *storage.Account) (int64, int64, error) {
	provider := GetProviderForAccount(account)
	if provider == nil {
		return 0, 0, fmt.Errorf("provider %s not supported", account.Provider)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return provider.FetchQuota(ctx, account)
}

