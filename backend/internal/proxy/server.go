package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"quotio-electron-go/backend/internal/providers"
	"quotio-electron-go/backend/internal/quota"
	"quotio-electron-go/backend/internal/storage"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Server struct {
	db              *gorm.DB
	port            int
	routingStrategy string
	proxy           *httputil.ReverseProxy
	server          *http.Server
	running         bool
	mu              sync.RWMutex
	router          *Router
	quotaTracker    *quota.Tracker
}

func NewServer(db *gorm.DB, port int, routingStrategy string) *Server {
	return &Server{
		db:              db,
		port:            port,
		routingStrategy: routingStrategy,
		router:          NewRouter(db, routingStrategy),
		quotaTracker:    quota.NewTracker(db),
	}
}

func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("proxy server already running")
	}

	// Create reverse proxy
	director := func(req *http.Request) {
		// Route request to appropriate provider (with validation)
		account, err := s.router.SelectAccount()
		if err != nil {
			log.Printf("Error selecting account: %v", err)
			req.URL.Scheme = "http"
			req.URL.Host = "localhost"
			req.URL.Path = "/error"
			return
		}

		// Pre-request validation: Check if account is valid for routing
		if !s.isAccountValidForRouting(account) {
			log.Printf("Account %d not valid for routing (status: %s)", account.ID, account.Status)
			// Try to select another account
			account, err = s.router.SelectNextAccount(account)
			if err != nil {
				log.Printf("No valid accounts available")
				req.URL.Scheme = "http"
				req.URL.Host = "localhost"
				req.URL.Path = "/error"
				return
			}
		}

		// Get provider
		provider := providers.GetProviderForAccount(account)
		if provider == nil {
			log.Printf("Provider not found: %s", account.Provider)
			req.URL.Scheme = "http"
			req.URL.Host = "localhost"
			req.URL.Path = "/error"
			return
		}

		// Get provider endpoint
		target, err := url.Parse(provider.GetBaseURL())
		if err != nil {
			log.Printf("Error parsing provider URL: %v", err)
			req.URL.Scheme = "http"
			req.URL.Host = "localhost"
			req.URL.Path = "/error"
			return
		}

		// Authenticate request using provider BEFORE modifying URL
		if err := provider.AuthenticateRequest(req, account); err != nil {
			log.Printf("Error authenticating request: %v", err)
			req.URL.Scheme = "http"
			req.URL.Host = "localhost"
			req.URL.Path = "/error"
			return
		}

		// Only set URL after successful authentication
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		// Add account ID to context for quota tracking
		req = req.WithContext(context.WithValue(req.Context(), "accountID", account.ID))
	}

	// Modify response to track quota and rate limits
	modifyResponse := func(resp *http.Response) error {
		accountID, ok := resp.Request.Context().Value("accountID").(uint)
		if !ok {
			return nil
		}

		// Determine success based on status code
		statusCode := resp.StatusCode
		success := statusCode >= 200 && statusCode < 300

		// Get account to find provider
		var account storage.Account
		if err := s.db.First(&account, accountID).Error; err == nil {
			provider := providers.GetProviderForAccount(&account)
			if provider != nil {
				// For non-streaming responses with JSON bodies, try to parse quota
				// For streaming responses (SSE, etc.), tokens are counted from stream size
				tokensUsed := int64(0)

				// Only buffer body for non-streaming, non-empty responses
				if resp.Header.Get("Content-Type") != "" &&
					!isStreamingResponse(resp) &&
					resp.ContentLength > 0 &&
					resp.ContentLength < 1*1024*1024 { // Only buffer responses < 1MB

					// Use TeeReader to buffer body without breaking streaming
					bufferedBody, newBody, err := teeResponseBody(resp)
					if err == nil && bufferedBody != nil {
						// Parse quota from buffered body
						tokensUsed, _ = provider.ParseQuotaFromBody(bufferedBody)
						// Replace body with new reader
						resp.Body = newBody
					}
				}

				// Also try header-based quota parsing (fallback if body parsing didn't work)
				if tokensUsed == 0 {
					tokensUsed, _ = provider.ParseQuotaFromResponse(resp)
				}

				// Parse rate limit headers using provider-specific config
				// Handle rate limits from headers
				client, _ := providers.CreateProviderClient(&account)
				if client != nil {
					rateLimits, err := client.ParseRateLimitsFromResponse(resp)
					if err == nil && rateLimits != nil {
						// Update account with auto-detected limits
						s.updateAccountRateLimits(accountID, rateLimits)

						// Check if we need to enter cooldown based on remaining quota
						if (rateLimits.TokensLimit > 0 && rateLimits.TokensRemaining == 0) ||
							(rateLimits.RequestsLimit > 0 && rateLimits.RequestsRemaining == 0) {
							log.Printf("Rate limit exhausted (headers) for account %d - entering cooldown", accountID)

							// Use reset time if available, otherwise default 15 mins
							cooldownDuration := 15 * time.Minute
							if !rateLimits.TokensReset.IsZero() {
								storage.SetAccountCooldown(accountID, rateLimits.TokensReset)
							} else if !rateLimits.RequestsReset.IsZero() {
								storage.SetAccountCooldown(accountID, rateLimits.RequestsReset)
							} else {
								storage.SetAccountCooldown(accountID, time.Now().Add(cooldownDuration))
							}
						}
					}
				}

				// Record usage
				s.quotaTracker.RecordUsage(accountID, tokensUsed, 1, statusCode, success)

				// Handle auth failures - increment consecutive failures before disabling
				if statusCode == 401 || statusCode == 403 {
					s.handleAuthFailure(accountID, resp)
				}

				// Check for rate limit from status code (429)
				if statusCode == 429 || provider.DetectRateLimit(resp) {
					log.Printf("Rate limit detected (status 429) for account %d", accountID)

					// If we didn't get a reset time from headers, set a default cooldown
					var accountCheck storage.Account
					if err := s.db.First(&accountCheck, accountID).Error; err == nil {
						if accountCheck.Status != "cooldown" {
							// 15 minute cooldown for 429 if no header info
							storage.SetAccountCooldown(accountID, time.Now().Add(15*time.Minute))
						}
					}
				}
			}
		}

		return nil
	}

	s.proxy = &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponse,
	}

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	s.running = true

	// Start server in goroutine
	go func() {
		log.Printf("Proxy server starting on port %d", s.port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Proxy server error: %v", err)
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
		}
	}()

	return nil
}

func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down proxy server: %v", err)
	}

	s.running = false
	log.Println("Proxy server stopped")
}

func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	if s.proxy == nil {
		http.Error(w, "Proxy not initialized", http.StatusInternalServerError)
		return
	}

	// Enforce API key if configured
	var proxyConfig storage.ProxyConfig
	if err := s.db.First(&proxyConfig).Error; err == nil && proxyConfig.APIKey != "" {
		auth := r.Header.Get("Authorization")
		const prefix = "Bearer "

		if !strings.HasPrefix(auth, prefix) || strings.TrimPrefix(auth, prefix) != proxyConfig.APIKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	s.proxy.ServeHTTP(w, r)
}

// isAccountValidForRouting checks if account is valid for routing
func (s *Server) isAccountValidForRouting(account *storage.Account) bool {
	// Skip disabled accounts
	if account.Status == "disabled" {
		return false
	}

	// Skip accounts in cooldown that haven't passed their reset time
	if account.Status == "cooldown" {
		if time.Now().Before(account.CooldownUntil) {
			return false
		}
		// Reactivate account if cooldown has passed
		storage.SetAccountStatus(account.ID, "active")
		account.Status = "active"
	}

	// Skip rate-limited accounts in fill_first mode
	if account.Status == "rate_limited" && s.routingStrategy == "fill_first" {
		return false
	}

	return true
}

// updateAccountRateLimits updates account with auto-detected rate limits
func (s *Server) updateAccountRateLimits(accountID uint, limits *providers.RateLimitInfo) {
	var account storage.Account
	if err := s.db.First(&account, accountID).Error; err != nil {
		return
	}

	// Only update quota limit if user hasn't manually set one
	if !account.QuotaManual {
		if limits.TokensLimit != 0 {
			account.QuotaLimit = limits.TokensLimit
			account.QuotaAutoDetected = true
		} else if limits.RequestsLimit != 0 {
			account.QuotaLimit = limits.RequestsLimit
			account.QuotaAutoDetected = true
		}
	}

	// Update rate limit remaining values
	account.RateLimitRequests = limits.RequestsLimit
	account.RateLimitRequestsRemaining = limits.RequestsRemaining
	account.RateLimitRequestsReset = limits.RequestsReset
	account.RateLimitTokens = limits.TokensLimit
	account.RateLimitTokensRemaining = limits.TokensRemaining
	account.RateLimitTokensReset = limits.TokensReset

	// Update used quota based on remaining (if not manually set)
	if !account.QuotaManual && limits.TokensLimit > 0 && limits.TokensRemaining >= 0 {
		account.QuotaUsed = limits.TokensLimit - limits.TokensRemaining
	}

	s.db.Save(&account)
}

// handleAuthFailure implements a smart auth failure strategy with retries
// It tracks consecutive auth failures and only disables after threshold
func (s *Server) handleAuthFailure(accountID uint, resp *http.Response) {
	var health storage.ProviderHealth

	// Get or create health record
	if err := s.db.Where("account_id = ?", accountID).FirstOrCreate(&health, storage.ProviderHealth{
		AccountID: accountID,
		IsHealthy: true,
	}).Error; err != nil {
		log.Printf("Error retrieving health for account %d: %v", accountID, err)
		return
	}

	// Check response body for specific error codes that indicate permanent auth issues
	if isPermanentAuthError(resp) {
		log.Printf("Permanent auth failure for account %d - marking disabled", accountID)
		storage.SetAccountStatus(accountID, "disabled")
		health.ConsecutiveFailures = 100 // Mark as permanently failed
		s.db.Save(&health)
		return
	}

	// Increment consecutive failures for transient issues
	health.ConsecutiveFailures++
	health.LastChecked = time.Now()

	const authFailureThreshold = 3 // Disable after 3 consecutive auth failures

	if health.ConsecutiveFailures >= authFailureThreshold {
		log.Printf("Auth failure threshold reached for account %d (consecutive failures: %d) - marking disabled",
			accountID, health.ConsecutiveFailures)
		storage.SetAccountStatus(accountID, "disabled")
	} else {
		log.Printf("Auth failure for account %d - warning state (consecutive failures: %d/%d)",
			accountID, health.ConsecutiveFailures, authFailureThreshold)
	}

	s.db.Save(&health)
}

// isPermanentAuthError checks response body for specific error codes indicating permanent auth issues
func isPermanentAuthError(resp *http.Response) bool {
	if resp == nil || resp.Body == nil {
		return false
	}

	// Read response body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// Reconstruct body for potential re-reading
	resp.Body = io.NopCloser(strings.NewReader(string(body)))

	respBody := string(body)

	// Check for permanent auth error indicators
	permanentErrors := []string{
		"invalid_api_key",        // Invalid API key
		"invalid_client_id",      // Invalid OAuth client
		"invalid_grant",          // Invalid OAuth grant (revoked token)
		"access_denied",          // User denied access
		"token_revoked",          // Token has been revoked
		"unauthorized_client",    // Client is unauthorized
		"unsupported_grant_type", // Unsupported grant type
	}

	for _, errCode := range permanentErrors {
		if strings.Contains(strings.ToLower(respBody), errCode) {
			return true
		}
	}

	return false
}

// isStreamingResponse detects if response is a streaming response (SSE, chunked, etc)
func isStreamingResponse(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	transferEncoding := resp.Header.Get("Transfer-Encoding")

	// Check for streaming content types
	streamingTypes := []string{
		"text/event-stream",        // Server-Sent Events
		"text/plain",               // Plain text streaming
		"application/octet-stream", // Binary streaming
	}

	for _, st := range streamingTypes {
		if strings.Contains(contentType, st) {
			return true
		}
	}

	// Check for chunked encoding
	return strings.Contains(strings.ToLower(transferEncoding), "chunked")
}

// teeResponseBody reads response body into a buffer while preserving it for the client
// Returns (bufferedData, newReadCloser, error)
func teeResponseBody(resp *http.Response) ([]byte, io.ReadCloser, error) {
	if resp.Body == nil {
		return nil, nil, nil
	}

	// Read entire body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Close original body
	resp.Body.Close()

	// Return buffered data and a new reader for the body
	return body, io.NopCloser(io.Reader(strings.NewReader(string(body)))), nil
}
