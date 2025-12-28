package notifications

import (
	"fmt"
	"log"
	"quotio-electron-go/backend/internal/storage"
	"time"
)

type Notifier struct {
	subscribers []NotificationSubscriber
}

type NotificationSubscriber interface {
	Notify(event NotificationEvent)
}

type NotificationEvent struct {
	Type      string    `json:"type"` // low_quota, rate_limit, cooldown, service_issue
	AccountID uint      `json:"account_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewNotifier() *Notifier {
	return &Notifier{
		subscribers: make([]NotificationSubscriber, 0),
	}
}

func (n *Notifier) Subscribe(subscriber NotificationSubscriber) {
	n.subscribers = append(n.subscribers, subscriber)
}

func (n *Notifier) NotifyLowQuota(account *storage.Account) {
	percentage := float64(account.QuotaUsed) / float64(account.QuotaLimit) * 100
	if percentage >= 80 {
		event := NotificationEvent{
			Type:      "low_quota",
			AccountID: account.ID,
			Message:   fmt.Sprintf("Account %s is at %.1f%% quota usage", account.Name, percentage),
			Timestamp: time.Now(),
		}
		n.broadcast(event)
	}
}

func (n *Notifier) NotifyRateLimit(account *storage.Account) {
	event := NotificationEvent{
		Type:      "rate_limit",
		AccountID: account.ID,
		Message:   fmt.Sprintf("Account %s has hit rate limit", account.Name),
		Timestamp: time.Now(),
	}
	n.broadcast(event)
}

func (n *Notifier) NotifyCooldown(account *storage.Account) {
	event := NotificationEvent{
		Type:      "cooldown",
		AccountID: account.ID,
		Message:   fmt.Sprintf("Account %s is in cooldown period", account.Name),
		Timestamp: time.Now(),
	}
	n.broadcast(event)
}

func (n *Notifier) NotifyServiceIssue(message string) {
	event := NotificationEvent{
		Type:      "service_issue",
		AccountID: 0,
		Message:   message,
		Timestamp: time.Now(),
	}
	n.broadcast(event)
}

func (n *Notifier) broadcast(event NotificationEvent) {
	log.Printf("Notification: %s - %s", event.Type, event.Message)
	for _, subscriber := range n.subscribers {
		go subscriber.Notify(event)
	}
}

