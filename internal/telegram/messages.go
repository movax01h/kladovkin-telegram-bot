package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// sendErrorMessage sends an error message to the user.
func (b *Bot) sendErrorMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

// SendNotification sends a notification to the user.
func (b *Bot) SendNotification(userID int64, text string) error {
	msg := tgbotapi.NewMessage(userID, text)
	_, err := b.api.Send(msg)
	return err
}

// SendErrorNotification sends an error notification to the admin.
func (b *Bot) SendErrorNotification(text string) {
	msg := tgbotapi.NewMessage(b.cfg.AdminID, text)
	b.api.Send(msg)
}
