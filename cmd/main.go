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

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
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
	logOutput, cleanup, err := tools.OpenLogFile(cfg.LoggerConfig.FilePath)
	if err != nil {
		log.Fatalf("failed to initialize log output: %s", err)
	}
	defer cleanup()
	logger.Setup(cfg.LoggerConfig.Level, logOutput)
	slog.Info("Logger is set up")

	// Initialize the database
	db, err := initializeDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Instantiate repositories
	userRepo := sqlite.NewSQLiteUserRepository(db)
	unitRepo := sqlite.NewSQLiteUnitRepository(db)
	subscriptionRepo := sqlite.NewSQLiteSubscriptionRepository(db)

	// Initialize the Telegram bot, passing in the repositories
	bot, err := telegram.NewBot(cfg.TelegramConfig, userRepo, subscriptionRepo)
	if err != nil {
		log.Fatalf("failed to initialize Telegram bot: %v", err)
	}
	slog.Info("Telegram bot initialized")

	// Initialize the notification service
	notificationService := notifier.NewNotifier(&cfg.NotifierConfig, userRepo, subscriptionRepo, bot)
	slog.Info("Notifier service initialized")

	// Initialize the parser
	parserService := parser.NewParser(&cfg.ParserConfig, userRepo, unitRepo, subscriptionRepo)
	slog.Info("Parser service initialized")

	// Initialize the scheduler
	schedulerService := scheduler.NewScheduler(&cfg.SchedulerConfig, subscriptionRepo, unitRepo, userRepo)
	slog.Info("")

	// Create a context and wait group for goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start the necessary goroutines (parser, notifier, bot, scheduler)
	startAllRoutines(ctx, &wg, bot, notificationService, parserService, schedulerService)

	// Wait for shutdown signal
	waitForShutdown(cancel, &wg)

	slog.Info("Shutting down application")
}

func initializeDatabase(cfg *config.Config) (*sql.DB, error) {
	// Ensure the data directory exists
	if err := ensureDataDirectoryExists(cfg.DatabaseConfig.Path); err != nil {
		return nil, err
	}

	// Initialize the SQLite database connection
	db, err := sqlite.NewSQLiteDB(cfg.DatabaseConfig.Path)
	if err != nil {
		return nil, err
	}

	// Ensure that the necessary tables are created
	if err := sqlite.InitializeDatabase(db); err != nil {
		return nil, err
	}

	slog.Info("SQLite database initialized")
	return db, nil
}

func ensureDataDirectoryExists(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func startAllRoutines(ctx context.Context, wg *sync.WaitGroup, b *telegram.Bot, n *notifier.Notifier, p *parser.Parser, s *scheduler.Scheduler) {
	startRoutine(ctx, wg, p.Start, "Error in parsing", b)
	startRoutine(ctx, wg, n.Start, "Error in notification routine", b)
	startRoutine(ctx, wg, b.Start, "Error in Telegram bot interaction", b)
	startRoutine(ctx, wg, s.ScheduleTasks, "Error in scheduling tasks", b)
}

func startRoutine(ctx context.Context, wg *sync.WaitGroup, routine func(ctx context.Context) error, errMsg string, bot *telegram.Bot) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := routine(ctx); err != nil {
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
