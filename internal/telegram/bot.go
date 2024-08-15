package telegram

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/movax01h/kladovkin-telegram-bot/config"
	"github.com/movax01h/kladovkin-telegram-bot/internal/repository"
)

// Bot represents the Telegram bot.Ø
type Bot struct {
	api              *tgbotapi.BotAPI
	userRepo         repository.UserRepository
	subscriptionRepo repository.SubscriptionRepository
}

// NewBot creates a new Bot instance.
func NewBot(cfg config.TelegramConfig, userRepo repository.UserRepository, subscriptionRepo repository.SubscriptionRepository) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		api:              api,
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
	}

	return bot, nil
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

// handleMessage processes incoming messages and executes commands.
func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.handleStart(ctx, message)
	case "subscribe":
		b.handleSubscribe(ctx, message)
	case "unsubscribe":
		b.handleUnsubscribe(ctx, message)
	default:
		b.handleUnknownCommand(ctx, message)
	}
}

// handleStart welcomes the user and provides basic instructions.
func (b *Bot) handleStart(ctx context.Context, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Use /subscribe to start receiving notifications and /unsubscribe to stop.")
	b.api.Send(msg)
}

// handleSubscribe subscribes the user to notifications.
func (b *Bot) handleSubscribe(ctx context.Context, message *tgbotapi.Message) {
	user, err := b.userRepo.GetByTelegramID(ctx, message.Chat.ID)
	if err != nil {
		slog.Error("Failed to retrieve user", "error", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error subscribing to notifications. Please try again later.")
		b.api.Send(msg)
		return
	}

	if user == nil {
		user = &repository.User{
			TelegramID: message.Chat.ID,
			Username:   message.From.UserName,
			Subscribed: true,
		}
		err = b.userRepo.Create(ctx, user)
		if err != nil {
			slog.Error("Failed to create user", "error", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Error subscribing to notifications. Please try again later.")
			b.api.Send(msg)
			return
		}
	} else {
		user.Subscribed = true
		err = b.userRepo.Update(ctx, user)
		if err != nil {
			slog.Error("Failed to update user subscription", "error", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Error subscribing to notifications. Please try again later.")
			b.api.Send(msg)
			return
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "You have been subscribed to notifications.")
	b.api.Send(msg)
}

// handleUnsubscribe unsubscribes the user from notifications.
func (b *Bot) handleUnsubscribe(ctx context.Context, message *tgbotapi.Message) {
	user, err := b.userRepo.GetByTelegramID(ctx, message.Chat.ID)
	if err != nil {
		slog.Error("Failed to retrieve user", "error", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error unsubscribing from notifications. Please try again later.")
		b.api.Send(msg)
		return
	}

	if user == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You are not subscribed to notifications.")
		b.api.Send(msg)
		return
	}

	user.Subscribed = false
	err = b.userRepo.Update(ctx, user)
	if err != nil {
		slog.Error("Failed to update user subscription", "error", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Error unsubscribing from notifications. Please try again later.")
		b.api.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "You have been unsubscribed from notifications.")
	b.api.Send(msg)
}

// handleUnknownCommand handles unrecognized commands.
func (b *Bot) handleUnknownCommand(ctx context.Context, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Please use /subscribe or /unsubscribe.")
	b.api.Send(msg)
}

// SendNotification sends a notification to a specific user.
func (b *Bot) SendNotification(user *repository.User, text string) {
	msg := tgbotapi.NewMessage(user.TelegramID, text)
	b.api.Send(msg)
}

// SendErrorNotification sends an error notification to the admin or a specific user.
func (b *Bot) SendErrorNotification(text string) {
	adminID := int64(123456789) // Replace with your admin Telegram ID
	msg := tgbotapi.NewMessage(adminID, text)
	b.api.Send(msg)
}
