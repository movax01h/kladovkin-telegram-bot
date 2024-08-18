package telegram

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/movax01h/kladovkin-telegram-bot/config"
	r "github.com/movax01h/kladovkin-telegram-bot/internal/repository"
)

type Bot struct {
	cfg              config.TelegramConfig
	api              *tgbotapi.BotAPI
	userRepo         r.UserRepository
	unitRepo         r.UnitRepository
	subscriptionRepo r.SubscriptionRepository
}

// NewBot creates a new Bot instance.
func NewBot(cfg config.TelegramConfig, userRepo r.UserRepository, unitRepo r.UnitRepository, subscriptionRepo r.SubscriptionRepository) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	return &Bot{
		cfg:              cfg,
		api:              api,
		userRepo:         userRepo,
		unitRepo:         unitRepo,
		subscriptionRepo: subscriptionRepo,
	}, nil
}

// Start begins polling for updates and handling messages.
func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				b.handleMessage(ctx, update.Message)
			}
		case <-ctx.Done():
			slog.Info("Telegram bot is shutting down")
			return ctx.Err()
		}
	}
}
