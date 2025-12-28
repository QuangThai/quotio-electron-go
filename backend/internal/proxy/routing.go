package proxy

import (
	"errors"
	"quotio-electron-go/backend/internal/storage"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
)

type Router struct {
	db              *gorm.DB
	strategy        string
	roundRobinIndex uint64
}

func NewRouter(db *gorm.DB, strategy string) *Router {
	return &Router{
		db:       db,
		strategy: strategy,
	}
}

func (r *Router) SelectAccount() (*storage.Account, error) {
	var accounts []storage.Account
	now := time.Now()

	// Include active accounts and cooldown accounts that have passed their reset time
	query := r.db.Where(
		"status = ? OR (status = ? AND cooldown_until < ?)",
		"active", "cooldown", now,
	)

	if err := query.Find(&accounts).Error; err != nil {
		return nil, err
	}

	// Reactivate accounts that have passed cooldown
	for i := range accounts {
		if accounts[i].Status == "cooldown" && time.Now().After(accounts[i].CooldownUntil) {
			r.db.Model(&accounts[i]).Update("status", "active")
			accounts[i].Status = "active"
		}
	}

	if len(accounts) == 0 {
		return nil, errors.New("no active accounts available")
	}

	switch r.strategy {
	case "round_robin":
		return r.selectRoundRobin(accounts)
	case "fill_first":
		return r.selectFillFirst(accounts)
	default:
		return r.selectRoundRobin(accounts)
	}
}

// SelectNextAccount tries to select the next valid account
func (r *Router) SelectNextAccount(excludeAccount *storage.Account) (*storage.Account, error) {
	var accounts []storage.Account
	now := time.Now()

	// Include active accounts and cooldown accounts that have passed their reset time
	// Exclude the current account to find a backup
	query := r.db.Where(
		"(status = ? OR (status = ? AND cooldown_until < ?)) AND id != ?",
		"active", "cooldown", now, excludeAccount.ID,
	)

	if err := query.Find(&accounts).Error; err != nil {
		return nil, err
	}

	// Reactivate accounts that have passed cooldown
	for i := range accounts {
		if accounts[i].Status == "cooldown" && time.Now().After(accounts[i].CooldownUntil) {
			r.db.Model(&accounts[i]).Update("status", "active")
			accounts[i].Status = "active"
		}
	}

	if len(accounts) == 0 {
		return nil, errors.New("no valid accounts available")
	}

	switch r.strategy {
	case "round_robin":
		return r.selectRoundRobin(accounts)
	case "fill_first":
		return r.selectFillFirst(accounts)
	default:
		return r.selectRoundRobin(accounts)
	}
}

func (r *Router) selectRoundRobin(accounts []storage.Account) (*storage.Account, error) {
	if len(accounts) == 0 {
		return nil, errors.New("no accounts available")
	}

	index := atomic.AddUint64(&r.roundRobinIndex, 1) - 1
	selected := accounts[index%uint64(len(accounts))]
	return &selected, nil
}

func (r *Router) selectFillFirst(accounts []storage.Account) (*storage.Account, error) {
	if len(accounts) == 0 {
		return nil, errors.New("no accounts available")
	}

	// Find account with quota remaining AND active status
	for i := range accounts {
		account := &accounts[i]

		// Check status (skip non-active)
		if account.Status != "active" {
			continue
		}

		// Check quota (allow 0 for unlimited)
		if account.QuotaLimit == 0 || account.QuotaUsed < account.QuotaLimit {
			return account, nil
		}
	}

	// All accounts exhausted or inactive - return first active one
	for i := range accounts {
		if accounts[i].Status == "active" {
			return &accounts[i], nil
		}
	}

	return &accounts[0], nil
}

// RefreshAccountStatus refreshes account status from database
func (r *Router) RefreshAccountStatus(accountID uint) error {
	var account storage.Account
	if err := r.db.First(&account, accountID).Error; err != nil {
		return err
	}

	// Check if account should be reactivated (cooldown expired)
	if account.Status == "cooldown" {
		// Simple cooldown logic - in production, check actual cooldown time
		account.Status = "active"
		return r.db.Save(&account).Error
	}

	return nil
}
