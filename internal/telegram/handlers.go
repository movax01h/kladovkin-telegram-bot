package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	m "github.com/movax01h/kladovkin-telegram-bot/internal/models"
	"log/slog"
	"time"
)

func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	if message.IsCommand() {
		// Handle command messages (like /start)
		switch message.Command() {
		case "start":
			b.handleStart(ctx, message)
		default:
			b.handleUnknownCommand(ctx, message)
		}
	} else {
		// Handle text messages (including button clicks)
		switch message.Text {
		case "New Subscription":
			b.handleNewSubscription(ctx, message)
		case "List Subscriptions":
			b.handleListSubscriptions(ctx, message)
		case "Return":
			b.handleReturn(ctx, message)
		default:
			b.handleUnknownCommand(ctx, message)
		}
	}
}

func (b *Bot) handleStart(ctx context.Context, message *tgbotapi.Message) {
	// Retrieve the user from the database
	user, err := b.userRepo.GetByTelegramID(message.Chat.ID)
	if err != nil {
		slog.Error("Failed to retrieve user", "error", err)
		b.sendErrorMessage(message.Chat.ID, "Error retrieving user. Please try again later.")
		return
	}

	// Check if the user is not found
	if user == nil {
		// Create a new user if not found
		user = &m.User{
			TelegramID: message.Chat.ID,
			UserName:   message.Chat.UserName,
			FirstName:  message.Chat.FirstName,
			LastName:   message.Chat.LastName,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err = b.userRepo.CreateUser(user)
		if err != nil {
			slog.Error("Failed to create user", "error", err)
			b.sendErrorMessage(message.Chat.ID, "Error creating user. Please try again later.")
			return
		}
	}

	// Send the welcome message
	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! What would you like to do?")
	msg.ReplyMarkup = b.mainMenu()
	b.api.Send(msg)
}

func (b *Bot) handleNewSubscription(ctx context.Context, message *tgbotapi.Message) {
	// Query available cities from the database
	cities, err := b.unitRepo.GetCities()
	if err != nil {
		slog.Error("Failed to retrieve cities", "error", err)
		b.sendErrorMessage(message.Chat.ID, "Error retrieving cities. Please try again later.")
		return
	}

	// Send cities as reply keyboard
	msg := tgbotapi.NewMessage(message.Chat.ID, "Select a city:")
	msg.ReplyMarkup = b.citySelectionKeyboard(cities)
	b.api.Send(msg)
}

func (b *Bot) handleCitySelection(ctx context.Context, message *tgbotapi.Message) {
	// Retrieve storages based on the selected city
	storages, err := b.unitRepo.GetStoragesByCity(message.Text)
	if err != nil {
		slog.Error("Failed to retrieve storages", "error", err)
		b.sendErrorMessage(message.Chat.ID, "Error retrieving storages. Please try again later.")
		return
	}

	// Send storages as reply keyboard
	msg := tgbotapi.NewMessage(message.Chat.ID, "Select a storage:")
	msg.ReplyMarkup = b.storageSelectionKeyboard(storages)
	b.api.Send(msg)
}

func (b *Bot) handleStorageSelection(ctx context.Context, message *tgbotapi.Message) {
	// Retrieve unit sizes based on the selected storage
	unitSizes, err := b.unitRepo.GetUnitSizesByStorage(message.Text)
	if err != nil {
		slog.Error("Failed to retrieve unit sizes", "error", err)
		b.sendErrorMessage(message.Chat.ID, "Error retrieving unit sizes. Please try again later.")
		return
	}

	// Send unit sizes as reply keyboard
	msg := tgbotapi.NewMessage(message.Chat.ID, "Select a unit size:")
	msg.ReplyMarkup = b.unitSizeSelectionKeyboard(unitSizes)
	b.api.Send(msg)
}

func (b *Bot) handleUnitSizeSelection(ctx context.Context, message *tgbotapi.Message) {
	// Here you would save the subscription details
	// For example, save the user ID and unit size to the subscription table

	// Confirm the subscription
	msg := tgbotapi.NewMessage(message.Chat.ID, "You have been subscribed.")
	b.api.Send(msg)
}

func (b *Bot) handleListSubscriptions(ctx context.Context, message *tgbotapi.Message) {
	// Get the user ID from the message
	user, err := b.userRepo.GetByTelegramID(message.Chat.ID)
	if err != nil {
		slog.Error("Failed to retrieve user", "error", err)
		b.sendErrorMessage(message.Chat.ID, "Error retrieving user. Please try again later.")
		return
	}

	// Check if the user is not found
	if user == nil {
		slog.Error("User not found", "telegram_id", message.Chat.ID)
		b.sendErrorMessage(message.Chat.ID, "User not found. Please start the bot again.")
		return
	}

	// Retrieve the list of subscriptions for the user
	subscriptions, err := b.subscriptionRepo.GetSubscriptionsByUserID(user.ID)
	if err != nil {
		slog.Error("Failed to retrieve subscriptions", "error", err)
		b.sendErrorMessage(message.Chat.ID, "Error retrieving subscriptions. Please try again later.")
		return
	}

	// Create a response message listing all subscriptions
	var response string
	for _, subscription := range subscriptions {
		response += fmt.Sprintf("City: %s, Storage: %s, Unit: %s\n", subscription.City, subscription.Storage, subscription.UnitSize)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	b.api.Send(msg)
}

func (b *Bot) handleReturn(ctx context.Context, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Returning to the main menu.")
	msg.ReplyMarkup = b.mainMenu() // Show the main menu again
	b.api.Send(msg)
}

func (b *Bot) handleUnknownCommand(ctx context.Context, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Please use the menu.")
	b.api.Send(msg)
}

// mainMenu creates the main menu keyboard.
func (b *Bot) mainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("New Subscription"),
			tgbotapi.NewKeyboardButton("List Subscriptions"),
		),
	)
}

func (b *Bot) citySelectionKeyboard(cities []string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton

	// Add each city as a button in a row
	for _, city := range cities {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(city)))
	}

	// Add the "Return" button as the last row
	returnButton := tgbotapi.NewKeyboardButton("Return")
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(returnButton))

	return tgbotapi.NewReplyKeyboard(rows...)
}

func (b *Bot) storageSelectionKeyboard(storages []string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	for _, storage := range storages {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(storage)))
	}
	return tgbotapi.NewReplyKeyboard(rows...)
}

func (b *Bot) unitSizeSelectionKeyboard(unitSizes []string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	for _, size := range unitSizes {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(size)))
	}
	return tgbotapi.NewReplyKeyboard(rows...)
}
