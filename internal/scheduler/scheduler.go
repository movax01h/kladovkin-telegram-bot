package scheduler

import (
	"context"
	"time"

	"github.com/movax01h/kladovkin-telegram-bot/config"
	"github.com/movax01h/kladovkin-telegram-bot/internal/repository"
	"log/slog"
)

// Scheduler is responsible for managing scheduled tasks.
type Scheduler struct {
	cfg              *config.SchedulerConfig
	subscriptionRepo repository.SubscriptionRepository
	unitRepo         repository.UnitRepository
	userRepo         repository.UserRepository
}

// NewScheduler creates a new Scheduler instance.
func NewScheduler(cfg *config.SchedulerConfig, subscriptionRepo repository.SubscriptionRepository, unitRepo repository.UnitRepository, userRepo repository.UserRepository) *Scheduler {
	return &Scheduler{
		cfg:              cfg,
		subscriptionRepo: subscriptionRepo,
		unitRepo:         unitRepo,
		userRepo:         userRepo,
	}
}

// ScheduleTasks starts the scheduler to run periodic tasks.
func (s *Scheduler) ScheduleTasks(ctx context.Context) error {
	// Example of scheduling a task to run once every day at a specific time
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			slog.Info("Running scheduled tasks")
			err := s.runDailyTasks(ctx)
			if err != nil {
				slog.Error("Error running daily tasks", "error", err)
			}
		case <-ctx.Done():
			slog.Info("Scheduler is shutting down")
			return ctx.Err()
		}
	}
}

// runDailyTasks contains the logic for tasks that should be run daily.
func (s *Scheduler) runDailyTasks(ctx context.Context) error {
	// Example task: Notify users with active subscriptions
	err := s.notifyActiveSubscribers(ctx)
	if err != nil {
		return err
	}

	// Add other tasks here as needed

	return nil
}

// notifyActiveSubscribers sends notifications to users with active subscriptions.
func (s *Scheduler) notifyActiveSubscribers(ctx context.Context) error {
	subscriptions, err := s.subscriptionRepo.GetActiveSubscriptions(ctx)
	if err != nil {
		return err
	}

	for _, subscription := range subscriptions {
		user, err := s.userRepo.GetByID(ctx, subscription.UserID)
		if err != nil {
			slog.Error("Failed to retrieve user for notification", "user_id", subscription.UserID, "error", err)
			continue
		}

		// Here, you would add the code to send the notification to the user.
		// For example:
		// err = sendNotification(user)
		// if err != nil {
		//    slog.Error("Failed to send notification", "user_id", user.ID, "error", err)
		// }

		slog.Info("Notification sent", "user_id", user.ID)
	}

	return nil
}
