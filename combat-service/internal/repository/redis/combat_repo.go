package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/kozie/lookism-rpg/combat-service/internal/domain"
)

type combatRepo struct {
	client *redis.Client
}

func NewCombatRepository(client *redis.Client) domain.CombatRepository {
	return &combatRepo{client: client}
}

func (r *combatRepo) SaveState(ctx context.Context, state *domain.CombatState) error {
	key := fmt.Sprintf("combat:%s", state.CombatID)
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	// Battles expire after 1 hour of inactivity
	return r.client.Set(ctx, key, data, time.Hour).Err()
}

func (r *combatRepo) GetState(ctx context.Context, combatID string) (*domain.CombatState, error) {
	key := fmt.Sprintf("combat:%s", combatID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	var state domain.CombatState
	err = json.Unmarshal([]byte(val), &state)
	return &state, err
}

func (r *combatRepo) DeleteState(ctx context.Context, combatID string) error {
	key := fmt.Sprintf("combat:%s", combatID)
	return r.client.Del(ctx, key).Err()
}
