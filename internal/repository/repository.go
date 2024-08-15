package repository

import (
	m "github.com/movax01h/kladovkin-telegram-bot/internal/models"
)

// UserRepository defines the methods to interact with the user data.
type UserRepository interface {
	CreateUser(user *m.User) error
	GetUserByID(id int64) (*m.User, error)
	GetAllUsers() ([]*m.User, error)
	UpdateUser(user *m.User) error
	DeleteUser(id int64) error
}

// UnitRepository defines the methods to interact with the unit data.
type UnitRepository interface {
	CreateUnit(unit *m.Unit) error
	GetUnitByID(id int64) (*m.Unit, error)
	GetAllUnits() ([]*m.Unit, error)
	UpdateUnit(unit *m.Unit) error
	DeleteUnit(id int64) error
}

// SubscriptionRepository defines the methods to interact with the subscription data.
type SubscriptionRepository interface {
	CreateSubscription(subscription *m.Subscription) error
	GetSubscriptionByID(id int64) (*m.Subscription, error)
	GetAllSubscriptions() ([]*m.Subscription, error)
	GetActiveSubscriptions() ([]*m.Subscription, error)
	UpdateSubscription(subscription *m.Subscription) error
	DeleteSubscription(id int64) error
}
