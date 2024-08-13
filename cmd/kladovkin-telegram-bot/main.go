package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/movax01h/kladovkin-telegram-bot/config"
	"github.com/movax01h/kladovkin-telegram-bot/internal/notifier"
	"github.com/movax01h/kladovkin-telegram-bot/internal/parser"
	"github.com/movax01h/kladovkin-telegram-bot/internal/repository/sqlite"
	"github.com/movax01h/kladovkin-telegram-bot/internal/scheduler"
	"github.com/movax01h/kladovkin-telegram-bot/internal/telegram"
	"github.com/movax01h/kladovkin-telegram-bot/pkg/logger"
	"github.com/movax01h/kladovkin-telegram-bot/pkg/tools"
	"log/slog"
)

func main() {
	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}

	// Initialize logger
	initializeLogger(cfg)

	// Initialize SQLite database and repositories
	db, err := initializeDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize application components
	bot, err := initializeBot(cfg)
	if err != nil {
		log.Fatalf("failed to initialize Telegram bot: %v", err)
	}

	// Initialize context and wait group
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start application routines
	startRoutines(ctx, cfg, &wg, bot, db)

	// Wait for shutdown signal
	waitForShutdown(ctx, cancel, &wg)

	// Final cleanup
	slog.Info("Shutting down application")
}

// initializeLogger sets up the logger with the configuration.
func initializeLogger(cfg *config.Config) {
	logOutput, cleanup, err := tools.OpenLogFile(cfg.LogFilePath)
	if err != nil {
		log.Fatalf("failed to initialize log output: %s", err)
	}
	defer cleanup()

	logger.Setup(cfg.LogLevel, logOutput)
	slog.Info("Logger is set up")
}

// initializeDatabase sets up the SQLite database and repositories.
func initializeDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}
	slog.Info("Database initialized")
	return db, nil
}

// initializeBot sets up the Telegram bot.
func initializeBot(cfg *config.Config) (*telegram.Bot, error) {
	bot, err := telegram.NewBot(cfg.Telegram)
	if err != nil {
		return nil, err
	}
	slog.Info("Telegram bot initialized")
	return bot, nil
}

// startRoutines launches the required goroutines for parsing, notifications, etc.
func startRoutines(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup, bot *telegram.Bot, db *sql.DB) {
	userRepo := sqlite.NewSQLiteUserRepository(db)
	unitRepo := sqlite.NewSQLiteUnitRepository(db)
	subscriptionRepo := sqlite.NewSQLiteSubscriptionRepository(db)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := parser.Start(ctx, cfg, userRepo, unitRepo); err != nil {
			slog.Error("Error in parsing", "error", err)
			notifyError(bot, err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := notifier.Start(ctx, cfg, userRepo, subscriptionRepo); err != nil {
			slog.Error("Error in notification routine", "error", err)
			notifyError(bot, err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := bot.Start(ctx); err != nil {
			slog.Error("Error in Telegram bot interaction", "error", err)
			notifyError(bot, err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := scheduler.ScheduleTasks(ctx, cfg, bot); err != nil {
			slog.Error("Error in scheduling tasks", "error", err)
			notifyError(bot, err)
		}
	}()
}

// waitForShutdown listens for OS signals to gracefully shut down the application.
func waitForShutdown(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stopChan:
		slog.Info("Received shutdown signal")
		cancel() // Trigger cancellation in the context
	case <-ctx.Done():
		slog.Info("Context cancelled")
	}

	// Wait for all Goroutines to finish
	wg.Wait()
}

// notifyError sends error messages directly to a personal Telegram account.
func notifyError(bot *telegram.Bot, err error) {
	if bot != nil {
		bot.SendErrorNotification("An error occurred: " + err.Error())
	}
}
