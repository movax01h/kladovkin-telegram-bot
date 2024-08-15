package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"log/slog"

	"github.com/movax01h/kladovkin-telegram-bot/config"
	"github.com/movax01h/kladovkin-telegram-bot/internal/notifier"
	"github.com/movax01h/kladovkin-telegram-bot/internal/parser"
	"github.com/movax01h/kladovkin-telegram-bot/internal/repository/sqlite"
	"github.com/movax01h/kladovkin-telegram-bot/internal/scheduler"
	"github.com/movax01h/kladovkin-telegram-bot/internal/telegram"
	"github.com/movax01h/kladovkin-telegram-bot/pkg/logger"
	"github.com/movax01h/kladovkin-telegram-bot/pkg/tools"
)

func main() {
	cfg, err := initializeConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}

	initializeLogger(cfg)

	db, err := initializeDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	bot, err := initializeBot(cfg)
	if err != nil {
		log.Fatalf("failed to initialize Telegram bot: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	startAllRoutines(ctx, &wg, cfg, bot, db)
	waitForShutdown(cancel, &wg)
	slog.Info("Shutting down application")
}

func initializeConfig() (*config.Config, error) {
	return config.NewConfig()
}

func initializeLogger(cfg *config.Config) {
	logOutput, cleanup, err := tools.OpenLogFile(cfg.LoggerConfig.FilePath)
	if err != nil {
		log.Fatalf("failed to initialize log output: %s", err)
	}
	defer cleanup()

	logger.Setup(cfg.LoggerConfig.Level, logOutput)
	slog.Info("Logger is set up")
}

func initializeDatabase(cfg *config.Config) (*sql.DB, error) {
	// Ensure the data directory exists
	err := ensureDataDirectoryExists(cfg.DatabaseConfig.Path)
	if err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}

	// Initialize the database
	db, err := sql.Open("sqlite3", cfg.DatabaseConfig.Path)
	if err != nil {
		return nil, err
	}
	slog.Info("Database initialized")
	return db, nil
}

func ensureDataDirectoryExists(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func initializeBot(cfg *config.Config) (*telegram.Bot, error) {
	bot, err := telegram.NewBot(cfg.TelegramConfig, nil, nil) // TODO: pass userRepo and subscriptionRepo
	if err != nil {
		return nil, err
	}
	slog.Info("Telegram bot initialized")
	return bot, nil
}

func startAllRoutines(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, bot *telegram.Bot, db *sql.DB) {
	userRepo := sqlite.NewSQLiteUserRepository(db)
	unitRepo := sqlite.NewSQLiteUnitRepository(db)
	subscriptionRepo := sqlite.NewSQLiteSubscriptionRepository(db)

	startRoutine(ctx, wg, func() error {
		return parser.Start(ctx, cfg, userRepo, unitRepo)
	}, "Error in parsing", bot)

	startRoutine(ctx, wg, func() error {
		return notifier.Start(ctx, cfg, userRepo, subscriptionRepo)
	}, "Error in notification routine", bot)

	startRoutine(ctx, wg, func() error {
		return bot.Start(ctx)
	}, "Error in Telegram bot interaction", bot)

	startRoutine(ctx, wg, func() error {
		return scheduler.ScheduleTasks(ctx, cfg, bot)
	}, "Error in scheduling tasks", bot)
}

func startRoutine(ctx context.Context, wg *sync.WaitGroup, routine func() error, errMsg string, bot *telegram.Bot) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := routine(); err != nil {
			slog.Error(errMsg, "error", err)
			notifyError(bot, err)
		}
	}()
}

func waitForShutdown(cancel context.CancelFunc, wg *sync.WaitGroup) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stopChan:
		slog.Info("Received shutdown signal")
		cancel()
	case <-context.Background().Done():
		slog.Info("Context cancelled")
	}

	wg.Wait()
}

func notifyError(bot *telegram.Bot, err error) {
	if bot != nil {
		bot.SendErrorNotification("An error occurred: " + err.Error())
	}
}
