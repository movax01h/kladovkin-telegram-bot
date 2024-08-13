package repository

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id int64) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context) ([]*User, error)
}

type UnitRepository interface {
	CreateUnit(ctx context.Context, unit *Unit) error
	GetUnitByID(ctx context.Context, id int64) (*Unit, error)
	UpdateUnit(ctx context.Context, unit *Unit) error
	DeleteUnit(ctx context.Context, id int64) error
	ListUnits(ctx context.Context) ([]*Unit, error)
}

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub *Subscription) error
	GetSubscriptionByID(ctx context.Context, id int64) (*Subscription, error)
	UpdateSubscription(ctx context.Context, sub *Subscription) error
	DeleteSubscription(ctx context.Context, id int64) error
	ListSubscriptions(ctx context.Context) ([]*Subscription, error)
}
