package api

import (
	"net/http"
	"os"
	"quotio-electron-go/backend/internal/agents"
	"quotio-electron-go/backend/internal/providers"
	"quotio-electron-go/backend/internal/proxy"
	"quotio-electron-go/backend/internal/storage"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now(),
	})
}

func (s *Server) handleDashboard(c *gin.Context) {
	// Get active accounts count
	var activeCount int64
	s.db.Model(&storage.Account{}).Where("status = ?", "active").Count(&activeCount)

	// Get total requests today
	var todayRequests int64
	startOfDay := time.Now().Truncate(24 * time.Hour)
	s.db.Model(&storage.QuotaHistory{}).
		Where("timestamp >= ?", startOfDay).
		Count(&todayRequests)

	// Get total tokens used today
	var tokensToday int64
	s.db.Model(&storage.QuotaHistory{}).
		Where("timestamp >= ?", startOfDay).
		Select("COALESCE(SUM(tokens_used),0)").
		Row().Scan(&tokensToday)

	// Compute success rate today
	var totalToday, successToday int64
	s.db.Model(&storage.QuotaHistory{}).
		Where("timestamp >= ?", startOfDay).
		Count(&totalToday)
	s.db.Model(&storage.QuotaHistory{}).
		Where("timestamp >= ? AND success = ?", startOfDay, true).
		Count(&successToday)

	var successRate float64
	if totalToday > 0 {
		successRate = float64(successToday) / float64(totalToday)
	}

	// Compute provider-level counts for badges
	type ProviderSummary struct {
		Provider string `json:"provider"`
		Accounts int64  `json:"accounts"`
	}

	var perProvider []ProviderSummary
	s.db.Model(&storage.Account{}).
		Select("provider, COUNT(*) as accounts").
		Where("status != ?", "disabled").
		Group("provider").
		Scan(&perProvider)

	// Get proxy status
	proxyRunning := s.proxy != nil && s.proxy.IsRunning()

	c.JSON(http.StatusOK, gin.H{
		"server_status":   "running",
		"proxy_running":   proxyRunning,
		"active_accounts": activeCount,
		"requests_today":  todayRequests,
		"tokens_today":    tokensToday,
		"success_rate":    successRate,
		"providers":       perProvider,
		"uptime":          time.Since(time.Now()).String(), // Simplified
	})
}

func (s *Server) handleGetProviders(c *gin.Context) {
	var accounts []storage.Account
	if err := s.db.Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enrich with health status
	type ProviderWithHealth struct {
		storage.Account
		IsHealthy   bool   `json:"is_healthy"`
		ResponseTime int64  `json:"response_time_ms"`
		LastChecked string `json:"last_checked"`
	}

	result := make([]ProviderWithHealth, 0, len(accounts))
	for _, account := range accounts {
		// Get health status
		var health storage.ProviderHealth
		if err := s.db.Where("account_id = ?", account.ID).First(&health).Error; err == nil {
			result = append(result, ProviderWithHealth{
				Account:      account,
				IsHealthy:    health.IsHealthy,
				ResponseTime: health.ResponseTime,
				LastChecked:  health.LastChecked.Format("2006-01-02 15:04:05"),
			})
		} else {
			result = append(result, ProviderWithHealth{
				Account:      account,
				IsHealthy:    true, // Default to healthy if no health record
				ResponseTime: 0,
				LastChecked:  "never",
			})
		}
	}

	// Don't expose sensitive data
	for i := range result {
		result[i].APIKey = ""
		result[i].OAuthToken = ""
		result[i].RefreshToken = ""
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) handleGetProviderStatus(c *gin.Context) {
	type ProviderStatus struct {
		Provider        string `json:"provider"`
		ActiveAccounts  int64  `json:"active_accounts"`
		LimitedAccounts int64  `json:"limited_accounts"`
		CooldownAccounts int64  `json:"cooldown_accounts"`
		TotalAccounts   int64  `json:"total_accounts"`
	}

	var statuses []ProviderStatus
	
	// Get all providers
	var providers []string
	s.db.Model(&storage.Account{}).
		Distinct("provider").
		Pluck("provider", &providers)

	// For each provider, get status breakdown
	for _, provider := range providers {
		status := ProviderStatus{Provider: provider}
		
		s.db.Model(&storage.Account{}).
			Where("provider = ? AND status = ?", provider, "active").
			Count(&status.ActiveAccounts)
		
		s.db.Model(&storage.Account{}).
			Where("provider = ? AND status = ?", provider, "rate_limited").
			Count(&status.LimitedAccounts)
		
		s.db.Model(&storage.Account{}).
			Where("provider = ? AND status = ?", provider, "cooldown").
			Count(&status.CooldownAccounts)
		
		s.db.Model(&storage.Account{}).
			Where("provider = ?", provider).
			Count(&status.TotalAccounts)
		
		statuses = append(statuses, status)
	}

	c.JSON(http.StatusOK, statuses)
}

func (s *Server) handleAddProvider(c *gin.Context) {
	var account storage.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if account.Status == "" {
		account.Status = "active"
	}
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	// Save account (must have ID for async validation)
	if err := s.db.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Asynchronously validate credentials in background
	go func(accountID uint) {
		var acc storage.Account
		if err := s.db.First(&acc, accountID).Error; err != nil {
			return
		}

		result := providers.ValidateAccountCredentials(&acc)
		if !result.IsValid {
			// Mark as disabled if validation fails
			s.db.Model(&acc).Update("status", "disabled")
		}

		// Record health
		s.db.Where("account_id = ?", accountID).
			Assign(&storage.ProviderHealth{
				AccountID:    accountID,
				IsHealthy:    result.IsValid,
				LastChecked:  time.Now(),
			}).
			FirstOrCreate(&storage.ProviderHealth{})
	}(account.ID)

	// Return immediately with account (status not yet validated)
	account.APIKey = ""
	account.OAuthToken = ""
	c.JSON(http.StatusCreated, account)
}

func (s *Server) handleUpdateProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var account storage.Account
	if err := s.db.First(&account, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account.UpdatedAt = time.Now()
	if err := s.db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	account.APIKey = ""
	account.OAuthToken = ""
	c.JSON(http.StatusOK, account)
}

func (s *Server) handleDetectProviderAccounts(c *gin.Context) {
	detected := []storage.Account{}

	// Check for environment-based provider accounts
	envChecks := map[string]string{
		"OPENAI_API_KEY":        "openai",
		"ANTHROPIC_API_KEY":     "claude",
		"GOOGLE_API_KEY":        "gemini",
		"GITHUB_TOKEN":          "copilot",
		"QWEN_API_KEY":          "qwen",
		"VERTEX_API_KEY":        "vertex",
		"ANTIGRAVITY_API_KEY":   "antigravity",
	}

	for envVar, provider := range envChecks {
		if key := os.Getenv(envVar); key != "" {
			detected = append(detected, storage.Account{
				Provider:     provider,
				Name:         provider + " (detected)",
				APIKey:       key,
				Status:       "active",
				QuotaLimit:   0,
				AutoDetected: true, // ✅ NEW: Mark as auto-detected
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			})
		}
	}

	c.JSON(http.StatusOK, detected)
}

func (s *Server) handleDeleteProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := s.db.Delete(&storage.Account{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}

func (s *Server) handleGetQuota(c *gin.Context) {
	accounts, err := storage.GetQuotaStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enrich with model usage and health status
	type QuotaWithModelsAndHealth struct {
		storage.Account
		ModelUsage   map[string]int64 `json:"model_usage"`
		IsHealthy    bool             `json:"is_healthy"`
		ResponseTime int64            `json:"response_time_ms"`
		LastChecked  string           `json:"last_checked"`
		// NEW: Auto-detected limits
		AutoDetectedLimit int64 `json:"auto_detected_limit"`
		IsManualQuota    bool  `json:"is_manual_quota"`
	}

	result := make([]QuotaWithModelsAndHealth, 0, len(accounts))
	for _, account := range accounts {
		// Exclude disabled accounts (invalid credentials)
		if account.Status == "disabled" {
			continue
		}

		modelUsage, _ := storage.GetQuotaByModel(account.ID)

		// Get health status
		var health storage.ProviderHealth
		isHealthy := true
		responseTime := int64(0)
		lastChecked := "never"
		if err := s.db.Where("account_id = ?", account.ID).First(&health).Error; err == nil {
			isHealthy = health.IsHealthy
			responseTime = health.ResponseTime
			lastChecked = health.LastChecked.Format("2006-01-02 15:04:05")
		}

		// Get auto-detected limit if not manual
		autoDetectedLimit := int64(0)
		if !account.QuotaManual {
			autoDetectedLimit = account.RateLimitTokens
		}

		result = append(result, QuotaWithModelsAndHealth{
			Account:           account,
			ModelUsage:        modelUsage,
			IsHealthy:         isHealthy,
			ResponseTime:      responseTime,
			LastChecked:       lastChecked,
			AutoDetectedLimit: autoDetectedLimit,
			IsManualQuota:     account.QuotaManual,
		})
	}

	// Don't expose sensitive data
	for i := range result {
		result[i].APIKey = ""
		result[i].OAuthToken = ""
		result[i].RefreshToken = ""
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) handleGetAgents(c *gin.Context) {
	// Detect agents with improved detection
	detected, err := agents.DetectAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Merge detected agents with stored configs
	var storedAgents []storage.AgentConfig
	s.db.Find(&storedAgents)

	// Create a map of stored agents
	storedMap := make(map[string]*storage.AgentConfig)
	for i := range storedAgents {
		storedMap[storedAgents[i].AgentName] = &storedAgents[i]
	}

	// Merge detected with stored
	type AgentWithStatus struct {
		storage.AgentConfig
		Configured       bool   `json:"configured"`
		ConfigPathExists bool   `json:"config_path_exists"`
		ValidationError   string `json:"validation_error,omitempty"`
	}

	result := make([]AgentWithStatus, 0, len(detected))
	for _, agent := range detected {
		config := storage.AgentConfig{
			AgentName:      agent.Name,
			ConfigPath:     agent.ConfigPath,
			Installed:      agent.Installed,
			AutoConfigured: false,
		}

		var validationError string
		if stored, ok := storedMap[agent.Name]; ok {
			config.ConfigPath = stored.ConfigPath
			config.ProxyURL = stored.ProxyURL
			config.AutoConfigured = stored.AutoConfigured
			config.LastConfigured = stored.LastConfigured
		}

		// Validate config if it exists
		if agent.ConfigPathExists {
			valid, err := agents.ValidateAgentConfig(agent.ConfigPath)
			if !valid {
				validationError = err
			}
		}

		result = append(result, AgentWithStatus{
			AgentConfig:      config,
			Configured:       agent.Configured,
			ConfigPathExists: agent.ConfigPathExists,
			ValidationError: validationError,
		})
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) handleConfigureAgent(c *gin.Context) {
	var req struct {
		AgentName string `json:"agent_name"`
		ProxyURL  string `json:"proxy_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Configure the agent
	if err := agents.ConfigureAgent(req.AgentName, req.ProxyURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save configuration to database
	var agent storage.AgentConfig
	// Get the default config path for this agent type
	configPath := ""
	for _, a := range agents.KnownAgents {
		if a.Name == req.AgentName {
			configPath = a.ConfigPath
			break
		}
	}

	if err := s.db.Where("agent_name = ?", req.AgentName).
		Assign(storage.AgentConfig{
			ConfigPath:     configPath,
			ProxyURL:       req.ProxyURL,
			AutoConfigured: false,
			LastConfigured: time.Now(),
		}).
		FirstOrCreate(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agent)
}

func (s *Server) handleRefreshAgents(c *gin.Context) {
	// Refresh agent detection
	detected, err := agents.DetectAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Merge with stored configs
	var storedAgents []storage.AgentConfig
	s.db.Find(&storedAgents)

	storedMap := make(map[string]*storage.AgentConfig)
	for i := range storedAgents {
		storedMap[storedAgents[i].AgentName] = &storedAgents[i]
	}

	type AgentWithStatus struct {
		storage.AgentConfig
		Configured       bool   `json:"configured"`
		ConfigPathExists bool   `json:"config_path_exists"`
		ValidationError   string `json:"validation_error,omitempty"`
	}

	result := make([]AgentWithStatus, 0, len(detected))
	for _, agent := range detected {
		config := storage.AgentConfig{
			AgentName:      agent.Name,
			ConfigPath:     agent.ConfigPath,
			Installed:      agent.Installed,
			AutoConfigured: false,
		}

		var validationError string
		if stored, ok := storedMap[agent.Name]; ok {
			config.ConfigPath = stored.ConfigPath
			config.ProxyURL = stored.ProxyURL
			config.AutoConfigured = stored.AutoConfigured
			config.LastConfigured = stored.LastConfigured
		}

		if agent.ConfigPathExists {
			valid, err := agents.ValidateAgentConfig(agent.ConfigPath)
			if !valid {
				validationError = err
			}
		}

		result = append(result, AgentWithStatus{
			AgentConfig:      config,
			Configured:       agent.Configured,
			ConfigPathExists: agent.ConfigPathExists,
			ValidationError: validationError,
		})
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) handleStartProxy(c *gin.Context) {
	if s.proxy == nil {
		var proxyConfig storage.ProxyConfig
		if err := s.db.First(&proxyConfig).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Proxy config not found"})
			return
		}

		s.proxy = proxy.NewServer(s.db, proxyConfig.Port, proxyConfig.RoutingStrategy)
	}

	if err := s.proxy.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proxy started"})
}

func (s *Server) handleStopProxy(c *gin.Context) {
	if s.proxy == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Proxy not running"})
		return
	}

	s.proxy.Stop()
	c.JSON(http.StatusOK, gin.H{"message": "Proxy stopped"})
}

func (s *Server) handleProxyStatus(c *gin.Context) {
	running := s.proxy != nil && s.proxy.IsRunning()
	c.JSON(http.StatusOK, gin.H{
		"running": running,
		"port":    s.config.ProxyPort,
	})
}

func (s *Server) handleGetSettings(c *gin.Context) {
	var proxyConfig storage.ProxyConfig
	if err := s.db.First(&proxyConfig).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, proxyConfig)
}

func (s *Server) handleUpdateSettings(c *gin.Context) {
	var proxyConfig storage.ProxyConfig
	if err := s.db.First(&proxyConfig).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&proxyConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.Save(&proxyConfig).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the server's config cache
	s.config.ProxyPort = proxyConfig.Port

	c.JSON(http.StatusOK, proxyConfig)
}

func (s *Server) handleGetProviderHealth(c *gin.Context) {
	var healths []storage.ProviderHealth
	if err := s.db.Find(&healths).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type HealthWithAccount struct {
		storage.ProviderHealth
		ProviderName string `json:"provider_name"`
		AccountName  string `json:"account_name"`
	}

	result := make([]HealthWithAccount, 0, len(healths))
	for _, health := range healths {
		var account storage.Account
		if err := s.db.First(&account, health.AccountID).Error; err == nil {
			result = append(result, HealthWithAccount{
				ProviderHealth: health,
				ProviderName:   account.Provider,
				AccountName:    account.Name,
			})
		}
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) handleCheckProviderHealth(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var account storage.Account
	if err := s.db.First(&account, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Validate credentials with real provider call
	result := providers.ValidateAccountCredentials(&account)

	// If valid, also try to fetch actual quota
	if result.IsValid {
		used, limit, err := providers.FetchProviderQuota(&account)
		if err == nil && (used > 0 || limit > 0) {
			storage.UpdateAccountRateLimits(account.ID, 0, 0, time.Time{}, limit, limit-used, time.Time{})
		}
	}

	now := time.Now()
	health := storage.ProviderHealth{
		AccountID:           account.ID,
		IsHealthy:           result.IsValid,
		ResponseTime:        0, // Could measure if needed
		LastChecked:         now,
		ConsecutiveFailures: 0,
	}

	// Upsert health record - properly update or create
	var existingHealth storage.ProviderHealth
	if err := s.db.Where("account_id = ?", account.ID).First(&existingHealth).Error; err != nil {
		// Create new record
		s.db.Create(&health)
	} else {
		// Update existing record
		existingHealth.IsHealthy = result.IsValid
		existingHealth.ResponseTime = 0
		existingHealth.LastChecked = now
		if !result.IsValid {
			existingHealth.ConsecutiveFailures++
		} else {
			existingHealth.ConsecutiveFailures = 0
		}
		s.db.Save(&existingHealth)
	}

	// Update account status if credentials are invalid
	if !result.IsValid {
		s.db.Model(&account).Update("status", "disabled")
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id":           account.ID,
		"is_healthy":           result.IsValid,
		"validation_reason":    result.Reason,
		"error_message":        result.ErrorMsg,
		"last_checked":         now,
	})
}

func (s *Server) handleGetQuotaHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	limit := 100 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	history, err := storage.GetQuotaHistory(uint(id), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (s *Server) handleGetFailedRequests(c *gin.Context) {
	// PHASE C: Get failed requests (likely from invalid credentials)
	limit := 50 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	type FailedRequest struct {
		ID        uint      `json:"id"`
		AccountID uint      `json:"account_id"`
		Provider  string    `json:"provider"`
		Name      string    `json:"account_name"`
		Model     string    `json:"model"`
		StatusCode int      `json:"status_code"`
		TokensUsed int64    `json:"tokens_used"`
		Timestamp time.Time `json:"timestamp"`
	}

	var failed []FailedRequest
	
	// Fetch failed requests with account info
	s.db.Table("quota_history").
		Joins("JOIN accounts ON quota_history.account_id = accounts.id").
		Select(
			"quota_history.id, quota_history.account_id, accounts.provider, accounts.name, quota_history.model, quota_history.status_code, quota_history.tokens_used, quota_history.timestamp",
		).
		Where("quota_history.success = ?", false).
		Order("quota_history.timestamp DESC").
		Limit(limit).
		Scan(&failed)

	c.JSON(http.StatusOK, failed)
}

func (s *Server) handleResetQuota(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := storage.ResetQuota(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reset account status to active
	storage.SetAccountStatus(uint(id), "active")

	c.JSON(http.StatusOK, gin.H{"message": "Quota reset successfully"})
}

func (s *Server) handleGetModels(c *gin.Context) {
	providerID := c.Query("provider")

	// If provider specified, return only that provider's models
	if providerID != "" {
		models := providers.GetModelsByProvider(providerID)
		if len(models) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found or no models available"})
			return
		}
		c.JSON(http.StatusOK, models)
		return
	}

	// Group models by provider for backward compatibility if needed, 
	// or just return the whole list. Standardizing on grouped list.
	type ProviderModels struct {
		Provider string               `json:"provider"`
		Models   []providers.ModelInfo `json:"models"`
	}

	grouped := make(map[string][]providers.ModelInfo)
	for _, m := range providers.SupportedModels {
		grouped[m.Provider] = append(grouped[m.Provider], m)
	}

	var result []ProviderModels
	for p, ms := range grouped {
		result = append(result, ProviderModels{
			Provider: p,
			Models:   ms,
		})
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) handleUpdateRoutingStrategy(c *gin.Context) {
	var req struct {
		Strategy string `json:"strategy" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate strategy
	if req.Strategy != "round_robin" && req.Strategy != "fill_first" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid strategy. Use 'round_robin' or 'fill_first'"})
		return
	}

	var proxyConfig storage.ProxyConfig
	if err := s.db.First(&proxyConfig).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	proxyConfig.RoutingStrategy = req.Strategy
	if err := s.db.Save(&proxyConfig).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Routing strategy updated", "strategy": req.Strategy})
}

func (s *Server) handleGetRateLimits(c *gin.Context) {
	accounts, err := storage.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type RateLimitStatus struct {
		AccountID              uint      `json:"account_id"`
		Provider               string    `json:"provider"`
		Name                   string    `json:"name"`
		Status                 string    `json:"status"`
		RequestsLimit          int64     `json:"requests_limit"`
		RequestsRemaining      int64     `json:"requests_remaining"`
		RequestsReset          time.Time `json:"requests_reset"`
		TokensLimit            int64     `json:"tokens_limit"`
		TokensRemaining        int64     `json:"tokens_remaining"`
		TokensReset            time.Time `json:"tokens_reset"`
		QuotaLimit             int64     `json:"quota_limit"`
		QuotaUsed              int64     `json:"quota_used"`
		QuotaAutoDetected      bool      `json:"quota_auto_detected"`
		QuotaManual            bool      `json:"quota_manual"`
		CooldownUntil          time.Time `json:"cooldown_until,omitempty"`
		ResponseTime           int64     `json:"response_time_ms"`
		LastChecked            string    `json:"last_checked"`
	}

	result := make([]RateLimitStatus, 0)

	for _, account := range accounts {
		if account.Status == "disabled" {
			continue
		}

		// Get health status
		var health storage.ProviderHealth
		responseTime := int64(0)
		lastChecked := "never"

		if err := s.db.Where("account_id = ?", account.ID).First(&health).Error; err == nil {
			responseTime = health.ResponseTime
			lastChecked = health.LastChecked.Format("2006-01-02 15:04:05")
		}

		result = append(result, RateLimitStatus{
			AccountID:         account.ID,
			Provider:          account.Provider,
			Name:              account.Name,
			Status:            account.Status,
			RequestsLimit:     account.RateLimitRequests,
			RequestsRemaining: account.RateLimitRequestsRemaining,
			RequestsReset:     account.RateLimitRequestsReset,
			TokensLimit:       account.RateLimitTokens,
			TokensRemaining:   account.RateLimitTokensRemaining,
			TokensReset:       account.RateLimitTokensReset,
			QuotaLimit:        account.QuotaLimit,
			QuotaUsed:         account.QuotaUsed,
			QuotaAutoDetected: account.QuotaAutoDetected,
			QuotaManual:       account.QuotaManual,
			CooldownUntil:     account.CooldownUntil,
			ResponseTime:      responseTime,
			LastChecked:       lastChecked,
		})
	}

	c.JSON(http.StatusOK, result)
}

// handleDetectOAuthCredentials detects OAuth credentials from CLI auth files
func (s *Server) handleDetectOAuthCredentials(c *gin.Context) {
	providerName := c.Query("provider")

	if providerName == "" {
		// Detect all providers
		allCreds := providers.DetectAllCredentials()
		result := make(map[string]gin.H)

		for p, creds := range allCreds {
			result[p] = gin.H{
				"found":      true,
				"expires_at": creds.ExpiresAt,
				"token_type": creds.TokenType,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"providers": result,
			"count":     len(allCreds),
		})
		return
	}

	// Detect specific provider
	creds, err := providers.DetectProviderCredentials(providerName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"found": false,
			"error": err.Error(),
		})
		return
	}

	if creds == nil || creds.AccessToken == "" {
		c.JSON(http.StatusOK, gin.H{
			"found":   false,
			"message": "No OAuth credentials found for " + providerName,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"found":      true,
		"expires_at": creds.ExpiresAt,
		"token_type": creds.TokenType,
	})
}

// handleAddProviderFromOAuth adds a provider using detected OAuth credentials
func (s *Server) handleAddProviderFromOAuth(c *gin.Context) {
	var req struct {
		Provider string `json:"provider" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Detect OAuth credentials
	creds, err := providers.DetectProviderCredentials(req.Provider)
	if err != nil || creds == nil || creds.AccessToken == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No OAuth credentials found for " + req.Provider,
		})
		return
	}

	// Create new account with OAuth token
	account := storage.Account{
		Provider:       req.Provider,
		Name:           req.Provider + " (OAuth)",
		OAuthToken:     creds.AccessToken,
		RefreshToken:   creds.RefreshToken,
		TokenExpiresAt: creds.ExpiresAt,
		Status:         "active",
		AutoDetected:   true,
	}

	if err := s.db.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Validate credentials asynchronously
	go func() {
		result := providers.ValidateAccountCredentials(&account)
		if !result.IsValid {
			storage.SetAccountStatus(account.ID, "disabled")
		} else {
			// Record health
			health := storage.ProviderHealth{
				AccountID:   account.ID,
				IsHealthy:   true,
				LastChecked: time.Now(),
			}
			s.db.Create(&health)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Provider added from OAuth credentials",
		"account": account.ID,
	})
}
