package domain

import "context"

// UserRepository defines the contract for user persistence.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
}

// CardRepository defines the contract for card inventory persistence.
type CardRepository interface {
	Create(ctx context.Context, card *Card) error
	GetByID(ctx context.Context, id string) (*Card, error)
	GetByUserID(ctx context.Context, userID string) ([]*Card, error)
	Update(ctx context.Context, card *Card) error
	Delete(ctx context.Context, id string) error
}

// SquadRepository defines the contract for squad persistence.
type SquadRepository interface {
	Upsert(ctx context.Context, squad *Squad) error
	GetByUserID(ctx context.Context, userID string) (*Squad, error)
}
