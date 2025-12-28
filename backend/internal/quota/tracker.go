package quota

import (
	"quotio-electron-go/backend/internal/storage"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Tracker struct {
	db       *gorm.DB
	mu       sync.RWMutex
	counters map[uint]*AccountCounter
}

type AccountCounter struct {
	RequestsCount int64
	TokensUsed    int64
	LastRequest   time.Time
}

func NewTracker(db *gorm.DB) *Tracker {
	return &Tracker{
		db:       db,
		counters: make(map[uint]*AccountCounter),
	}
}

func (t *Tracker) RecordUsage(accountID uint, tokensUsed int64, requestsCount int, statusCode int, success bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Update in-memory counter
	if t.counters[accountID] == nil {
		t.counters[accountID] = &AccountCounter{}
	}
	counter := t.counters[accountID]
	counter.RequestsCount += int64(requestsCount)
	counter.TokensUsed += tokensUsed
	counter.LastRequest = time.Now()

	// Update database (async to avoid blocking)
	go func() {
		storage.UpdateQuotaUsage(accountID, tokensUsed, requestsCount)
		storage.RecordQuotaHistory(accountID, tokensUsed, requestsCount, statusCode, success)
	}()
}

func (t *Tracker) GetUsage(accountID uint) (int64, int64) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if counter, ok := t.counters[accountID]; ok {
		return counter.TokensUsed, counter.RequestsCount
	}
	return 0, 0
}

func (t *Tracker) CheckRateLimit(accountID uint) (bool, error) {
	var account storage.Account
	if err := t.db.First(&account, accountID).Error; err != nil {
		return false, err
	}

	// Check if account is rate limited
	if account.Status == "rate_limited" || account.Status == "cooldown" {
		return true, nil
	}

	// Check quota limit
	if account.QuotaLimit > 0 && account.QuotaUsed >= account.QuotaLimit {
		storage.SetAccountStatus(accountID, "rate_limited")
		return true, nil
	}

	return false, nil
}

func (t *Tracker) ResetAccount(accountID uint) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.counters, accountID)
	return storage.ResetQuota(accountID)
}

