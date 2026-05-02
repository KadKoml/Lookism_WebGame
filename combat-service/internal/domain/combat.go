package domain

import (
	pb "github.com/kozie/lookism-rpg/api/proto/combat"
)

// CombatState represents the full snapshot of an ongoing battle.
type CombatState struct {
	CombatID        string               `json:"combat_id"`
	PlayerID        string               `json:"player_id"`
	CurrentTurn     int32                `json:"current_turn"`
	PlayerTeam      []*pb.Character      `json:"player_team"`
	EnemyTeam       []*pb.Character      `json:"enemy_team"`
	PlayerTeamState []*pb.CharacterState `json:"player_team_state"`
	EnemyTeamState  []*pb.CharacterState `json:"enemy_team_state"`
	IsEnded         bool                 `json:"is_ended"`
	Winner          string               `json:"winner"`
}
