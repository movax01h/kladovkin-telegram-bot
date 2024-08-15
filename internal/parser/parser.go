package parser

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/movax01h/kladovkin-telegram-bot/config"
	"github.com/movax01h/kladovkin-telegram-bot/internal/repository"
	"golang.org/x/net/html"
)

// Parser handles the logic for parsing HTML data.
type Parser struct {
	cfg              *config.Config
	userRepo         repository.UserRepository
	unitRepo         repository.UnitRepository
	subscriptionRepo repository.SubscriptionRepository
	client           *http.Client
}

// NewParser creates a new Parser instance.
func NewParser(cfg *config.Config, userRepo repository.UserRepository, unitRepo repository.UnitRepository, subscriptionRepo repository.SubscriptionRepository) *Parser {
	return &Parser{
		cfg:              cfg,
		userRepo:         userRepo,
		unitRepo:         unitRepo,
		subscriptionRepo: subscriptionRepo,
		client:           &http.Client{Timeout: 10 * time.Second},
	}
}

// Start initiates the parsing process and runs it in a loop.
func (p *Parser) Start(ctx context.Context) error {
	slog.Info("Parser started")
	ticker := time.NewTicker(p.cfg.Parser.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Parser shutting down")
			return nil
		case <-ticker.C:
			slog.Info("Running parser")
			if err := p.parseAndStoreData(); err != nil {
				slog.Error("Failed to parse and store data", "error", err)
				return err
			}
		}
	}
}

// parseAndStoreData fetches the HTML data, parses it, and stores the relevant information.
func (p *Parser) parseAndStoreData() error {
	resp, err := p.client.Get(p.cfg.Parser.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	// Extract data from the parsed HTML
	users, units, subscriptions := p.extractData(doc)

	// Store the extracted data in the database
	if err := p.storeData(users, units, subscriptions); err != nil {
		return err
	}

	return nil
}

// extractData extracts the user, unit, and subscription information from the parsed HTML document.
func (p *Parser) extractData(doc *html.Node) ([]repository.User, []repository.Unit, []repository.Subscription) {
	var users []repository.User
	var units []repository.Unit
	var subscriptions []repository.Subscription

	// TODO: Implement the HTML parsing logic here to extract relevant data
	// You would navigate through the HTML nodes, looking for tables, rows, and columns,
	// and then populate the users, units, and subscriptions slices with the extracted data.

	return users, units, subscriptions
}

// storeData saves the extracted data into the database using the repositories.
func (p *Parser) storeData(users []repository.User, units []repository.Unit, subscriptions []repository.Subscription) error {
	// Store users
	for _, user := range users {
		if err := p.userRepo.SaveUser(user); err != nil {
			slog.Error("Failed to save user", "userID", user.ID, "error", err)
			continue
		}
	}

	// Store units
	for _, unit := range units {
		if err := p.unitRepo.SaveUnit(unit); err != nil {
			slog.Error("Failed to save unit", "unitID", unit.ID, "error", err)
			continue
		}
	}

	// Store subscriptions
	for _, subscription := range subscriptions {
		if err := p.subscriptionRepo.SaveSubscription(subscription); err != nil {
			slog.Error("Failed to save subscription", "subscriptionID", subscription.ID, "error", err)
			continue
		}
	}

	return nil
}
