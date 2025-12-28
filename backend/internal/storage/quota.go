package storage

import (
	"time"

	"gorm.io/gorm"
)

// GetQuotaStatus returns current quota status for all accounts
func GetQuotaStatus() ([]Account, error) {
	var accounts []Account
	err := DB.Where("status != ?", "disabled").Find(&accounts).Error
	return accounts, err
}

// UpdateQuotaUsage updates quota usage for an account
// Uses atomic SQL update to prevent race conditions on concurrent requests
func UpdateQuotaUsage(accountID uint, tokensUsed int64, requestsCount int) error {
	// Atomically increment quota_used and update last_used timestamp
	result := DB.Model(&Account{}).
		Where("id = ?", accountID).
		Updates(map[string]interface{}{
			"quota_used": gorm.Expr("quota_used + ?", tokensUsed),
			"last_used":  time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	// Now check if quota limit is reached and update status if needed
	var account Account
	if err := DB.First(&account, accountID).Error; err != nil {
		return err
	}

	// Update status to rate_limited if quota limit is exceeded
	if account.QuotaLimit > 0 && account.QuotaUsed >= account.QuotaLimit && account.Status != "rate_limited" {
		return DB.Model(&Account{}).Where("id = ?", accountID).Update("status", "rate_limited").Error
	}

	return nil
}

// RecordQuotaHistory records quota usage in history
func RecordQuotaHistory(accountID uint, tokensUsed int64, requestsCount int, statusCode int, success bool) error {
	history := QuotaHistory{
		AccountID:     accountID,
		TokensUsed:    tokensUsed,
		RequestsCount: requestsCount,
		StatusCode:    statusCode,
		Success:       success,
		Timestamp:     time.Now(),
	}
	return DB.Create(&history).Error
}

// RecordModelUsage records usage for a specific model
func RecordModelUsage(accountID uint, model string, tokensUsed int64) error {
	// Update quota
	if err := UpdateQuotaUsage(accountID, tokensUsed, 1); err != nil {
		return err
	}

	// Record history with model info
	history := QuotaHistory{
		AccountID:     accountID,
		TokensUsed:    tokensUsed,
		RequestsCount: 1,
		Model:         model,
		StatusCode:    200,
		Success:       true,
		Timestamp:     time.Now(),
	}
	return DB.Create(&history).Error
}

// GetQuotaByModel returns quota usage grouped by model for an account
func GetQuotaByModel(accountID uint) (map[string]int64, error) {
	var results []struct {
		Model      string
		TokensUsed int64
	}

	err := DB.Model(&QuotaHistory{}).
		Select("model, SUM(tokens_used) as tokens_used").
		Where("account_id = ?", accountID).
		Group("model").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to map
	modelQuota := make(map[string]int64)
	for _, r := range results {
		modelQuota[r.Model] = r.TokensUsed
	}

	return modelQuota, nil
}

// GetQuotaHistory returns quota history for an account
func GetQuotaHistory(accountID uint, limit int) ([]QuotaHistory, error) {
	var history []QuotaHistory
	query := DB.Where("account_id = ?", accountID).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&history).Error
	return history, err
}

// SetAccountStatus updates account status
func SetAccountStatus(accountID uint, status string) error {
	return DB.Model(&Account{}).Where("id = ?", accountID).Update("status", status).Error
}

// SetAccountCooldown sets account cooldown status and time
func SetAccountCooldown(accountID uint, cooldownUntil time.Time) error {
	return DB.Model(&Account{}).Where("id = ?", accountID).Updates(map[string]interface{}{
		"status":             "cooldown",
		"cooldown_until":     cooldownUntil,
		"last_rate_limit_at": time.Now(),
	}).Error
}

// ReactivateAccountFromCooldown reactivates an account if cooldown has passed
func ReactivateAccountFromCooldown(accountID uint) error {
	var account Account
	if err := DB.First(&account, accountID).Error; err != nil {
		return err
	}

	// Only reactivate if cooldown has passed
	if account.Status == "cooldown" && time.Now().After(account.CooldownUntil) {
		return SetAccountStatus(accountID, "active")
	}
	return nil
}

// ResetQuota resets quota usage for an account
func ResetQuota(accountID uint) error {
	return DB.Model(&Account{}).Where("id = ?", accountID).Updates(map[string]interface{}{
		"quota_used":     0,
		"status":         "active",
		"cooldown_until": time.Time{},
	}).Error
}

// GetAllAccounts returns all accounts (for rate limits endpoint)
func GetAllAccounts() ([]Account, error) {
	var accounts []Account
	err := DB.Find(&accounts).Error
	return accounts, err
}

// UpdateAccountRateLimits updates rate limit fields from provider headers
// Uses atomic SQL updates to prevent race conditions on concurrent requests
func UpdateAccountRateLimits(accountID uint, requestsLimit int64, requestsRemaining int64, requestsReset time.Time, tokensLimit int64, tokensRemaining int64, tokensReset time.Time) error {
	var account Account
	if err := DB.First(&account, accountID).Error; err != nil {
		return err
	}

	// Build updates map for atomic update
	updates := map[string]interface{}{
		"rate_limit_requests":           requestsLimit,
		"rate_limit_requests_remaining": requestsRemaining,
		"rate_limit_requests_reset":     requestsReset,
		"rate_limit_tokens":             tokensLimit,
		"rate_limit_tokens_remaining":   tokensRemaining,
		"rate_limit_tokens_reset":       tokensReset,
	}

	// Only update auto-detected quota limit (not user manual override)
	if !account.QuotaManual {
		if tokensLimit > 0 {
			updates["quota_limit"] = tokensLimit
			updates["quota_auto_detected"] = true
		}

		// Update used quota based on remaining (only if not manually set)
		if tokensLimit > 0 && tokensRemaining >= 0 {
			updates["quota_used"] = tokensLimit - tokensRemaining
		}
	}

	// Perform atomic update
	return DB.Model(&Account{}).Where("id = ?", accountID).Updates(updates).Error
}

// GetAccountsInCooldown returns accounts that are in cooldown but can be reactivated
func GetAccountsInCooldown() ([]Account, error) {
	var accounts []Account
	now := time.Now()
	err := DB.Where("status = ? AND cooldown_until < ?", "cooldown", now).Find(&accounts).Error
	return accounts, err
}

// GetAllFailedRequests returns all failed quota history entries
func GetAllFailedRequests(limit int) ([]QuotaHistory, error) {
	var history []QuotaHistory
	query := DB.Where("success = ?", false).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&history).Error
	return history, err
}
