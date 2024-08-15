package notifier

import (
	"context"
	m "github.com/movax01h/kladovkin-telegram-bot/internal/models"
	"log/slog"
	"time"

	"github.com/movax01h/kladovkin-telegram-bot/config"
	"github.com/movax01h/kladovkin-telegram-bot/internal/repository"
	"github.com/movax01h/kladovkin-telegram-bot/internal/telegram"
)

// Notifier handles the logic for sending notifications.
type Notifier struct {
	cfg              *config.NotifierConfig
	userRepo         repository.UserRepository
	subscriptionRepo repository.SubscriptionRepository
	telegramBot      *telegram.Bot
}

// NewNotifier creates a new Notifier instance.
func NewNotifier(cfg *config.NotifierConfig, userRepo repository.UserRepository, subscriptionRepo repository.SubscriptionRepository, telegramBot *telegram.Bot) *Notifier {
	return &Notifier{
		cfg:              cfg,
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
		telegramBot:      telegramBot,
	}
}

// Start begins the notification process.
func (n *Notifier) Start(ctx context.Context) error {
	slog.Info("Notifier started")
	ticker := time.NewTicker(24 * time.Hour) // Set the ticker to send notifications once a day
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Notifier shutting down")
			return nil
		case <-ticker.C:
			slog.Info("Sending daily notifications")
			if err := n.sendDailyNotifications(); err != nil {
				slog.Error("Failed to send daily notifications", "error", err)
				return err
			}
		}
	}
}

// sendDailyNotifications sends notifications to users who meet the criteria.
func (n *Notifier) sendDailyNotifications() error {
	// Fetch active subscriptions
	activeSubscriptions, err := n.subscriptionRepo.GetActiveSubscriptions()
	if err != nil {
		return err
	}

	// Iterate over active subscriptions and send notifications
	for _, subscription := range activeSubscriptions {
		// Check if the user meets the criteria for a notification
		user, err := n.userRepo.GetUserByID(subscription.UserID)
		if err != nil {
			slog.Error("Failed to retrieve user", "userID", subscription.UserID, "error", err)
			continue
		}

		// Avoid spamming by checking last notification timestamp
		if shouldNotify(user) {
			message := "Your subscription is active. Don't forget to check our updates!"
			err := n.telegramBot.SendNotification(user, message)
			if err != nil {
				slog.Error("Failed to send notification", "userID", user.ID, "error", err)
				continue
			}
			// Update last notification time
			user.LastNotified = time.Now()
			err = n.userRepo.UpdateUser(user)
			if err != nil {
				slog.Error("Failed to update user's last notification time", "userID", user.ID, "error", err)
			}
		}
	}

	return nil
}

// shouldNotify checks if the user should receive a notification based on the last notified timestamp.
func shouldNotify(user *m.User) bool {
	// Example logic: Notify if more than 24 hours have passed since the last notification
	return time.Since(user.LastNotified) > 24*time.Hour
}
