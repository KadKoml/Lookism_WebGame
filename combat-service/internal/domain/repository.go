package domain

import (
	"context"
)

// CombatRepository handles storing and retrieving active battle states.
// This is typically backed by Redis since battles are fast-paced and ephemeral.
type CombatRepository interface {
	SaveState(ctx context.Context, state *CombatState) error
	GetState(ctx context.Context, combatID string) (*CombatState, error)
	DeleteState(ctx context.Context, combatID string) error
}
