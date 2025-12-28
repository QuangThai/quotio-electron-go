package storage

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// Account represents an AI provider account
type Account struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	Provider       string    `gorm:"not null;index" json:"provider"`
	Name           string    `json:"name"`
	APIKey         string    `gorm:"type:text" json:"api_key,omitempty"`       // Encrypted in production
	OAuthToken     string    `gorm:"type:text" json:"oauth_token,omitempty"`   // Encrypted in production
	RefreshToken   string    `gorm:"type:text" json:"refresh_token,omitempty"` // OAuth refresh token
	TokenExpiresAt time.Time `json:"token_expires_at"`

	// Hybrid quota fields
	QuotaLimit        int64 `json:"quota_limit"` // 0 = unlimited or auto-detect
	QuotaUsed         int64 `json:"quota_used"`
	QuotaManual       bool  `gorm:"default:false" json:"quota_manual"`        // true if user manually set
	QuotaAutoDetected bool  `gorm:"default:false" json:"quota_auto_detected"` // true if from provider headers

	// Auto-detected rate limits (from provider headers)
	RateLimitRequests          int64     `json:"rate_limit_requests"`
	RateLimitRequestsRemaining int64     `json:"rate_limit_requests_remaining"`
	RateLimitRequestsReset     time.Time `json:"rate_limit_requests_reset"`
	RateLimitTokens            int64     `json:"rate_limit_tokens"`
	RateLimitTokensRemaining   int64     `json:"rate_limit_tokens_remaining"`
	RateLimitTokensReset       time.Time `json:"rate_limit_tokens_reset"`

	// Cooldown management
	CooldownUntil   time.Time `json:"cooldown_until"`
	LastRateLimitAt time.Time `json:"last_rate_limit_at"`

	Status             string    `gorm:"default:active" json:"status"`             // active, rate_limited, cooldown, disabled
	AutoDetected       bool      `gorm:"default:false" json:"auto_detected"`       // True if from env vars
	SupportsManualAuth bool      `gorm:"default:true" json:"supports_manual_auth"` // Can add manually
	ModelAccess        string    `gorm:"type:text" json:"model_access"`            // JSON array of models
	Priority           int       `gorm:"default:0" json:"priority"`                // For routing
	LastUsed           time.Time `json:"last_used"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// BeforeSave hook encrypts sensitive fields before saving to database
func (a *Account) BeforeSave(tx *gorm.DB) error {
	var err error

	// Encrypt APIKey if not empty
	if a.APIKey != "" {
		a.APIKey, err = Encrypt(a.APIKey)
		if err != nil {
			log.Printf("Error encrypting APIKey: %v", err)
			return err
		}
	}

	// Encrypt OAuthToken if not empty
	if a.OAuthToken != "" {
		a.OAuthToken, err = Encrypt(a.OAuthToken)
		if err != nil {
			log.Printf("Error encrypting OAuthToken: %v", err)
			return err
		}
	}

	// Encrypt RefreshToken if not empty
	if a.RefreshToken != "" {
		a.RefreshToken, err = Encrypt(a.RefreshToken)
		if err != nil {
			log.Printf("Error encrypting RefreshToken: %v", err)
			return err
		}
	}

	return nil
}

// AfterFind hook decrypts sensitive fields after loading from database
func (a *Account) AfterFind(tx *gorm.DB) error {
	var err error

	// Decrypt APIKey if not empty
	if a.APIKey != "" {
		a.APIKey, err = Decrypt(a.APIKey)
		if err != nil {
			log.Printf("Error decrypting APIKey: %v", err)
			return err
		}
	}

	// Decrypt OAuthToken if not empty
	if a.OAuthToken != "" {
		a.OAuthToken, err = Decrypt(a.OAuthToken)
		if err != nil {
			log.Printf("Error decrypting OAuthToken: %v", err)
			return err
		}
	}

	// Decrypt RefreshToken if not empty
	if a.RefreshToken != "" {
		a.RefreshToken, err = Decrypt(a.RefreshToken)
		if err != nil {
			log.Printf("Error decrypting RefreshToken: %v", err)
			return err
		}
	}

	return nil
}

// QuotaHistory tracks historical quota usage
type QuotaHistory struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	AccountID     uint      `gorm:"not null;index" json:"account_id"`
	Account       Account   `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	RequestsCount int       `json:"requests_count"`
	TokensUsed    int64     `json:"tokens_used"`
	Model         string    `json:"model"` // Model used (e.g., "claude-3-opus")
	StatusCode    int       `json:"status_code"`
	Success       bool      `json:"success"`
	Timestamp     time.Time `gorm:"index" json:"timestamp"`
}

// ProviderHealth tracks health status of provider accounts
type ProviderHealth struct {
	ID                  uint      `gorm:"primarykey" json:"id"`
	AccountID           uint      `gorm:"not null;uniqueIndex" json:"account_id"`
	Account             Account   `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	IsHealthy           bool      `gorm:"default:true" json:"is_healthy"`
	ResponseTime        int64     `json:"response_time_ms"` // Response time in milliseconds
	LastChecked         time.Time `json:"last_checked"`
	ConsecutiveFailures int       `gorm:"default:0" json:"consecutive_failures"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ProxyConfig stores proxy server configuration
type ProxyConfig struct {
	ID              uint   `gorm:"primarykey" json:"id"`
	Port            int    `gorm:"default:8081" json:"port"`
	RoutingStrategy string `gorm:"default:round_robin" json:"routing_strategy"` // round_robin, fill_first
	AutoStart       bool   `gorm:"default:false" json:"auto_start"`
	APIKey          string `gorm:"type:text" json:"api_key"` // API key for proxy authentication
}

// AgentConfig stores agent configuration
type AgentConfig struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	AgentName      string    `gorm:"uniqueIndex" json:"agent_name"`
	ConfigPath     string    `json:"config_path"`
	Installed      bool      `json:"installed"`
	AutoConfigured bool      `gorm:"default:false" json:"auto_configured"`
	ProxyURL       string    `json:"proxy_url"`
	LastConfigured time.Time `json:"last_configured"`
}
