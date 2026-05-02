package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	pb "github.com/kozie/lookism-rpg/api/proto/combat"
	userpb "github.com/kozie/lookism-rpg/api/proto/user"
	"github.com/kozie/lookism-rpg/combat-service/internal/domain"
	"google.golang.org/grpc"
)

var (
	ErrCombatNotFound = errors.New("combat not found or ended")
	ErrInvalidTarget  = errors.New("invalid target")
	ErrInvalidAction  = errors.New("invalid action")
)

type CombatUsecase struct {
	repo       domain.CombatRepository
	nc         *nats.Conn
	userClient userpb.UserServiceClient
}

func NewCombatUsecase(repo domain.CombatRepository, nc *nats.Conn, userConn *grpc.ClientConn) *CombatUsecase {
	return &CombatUsecase{
		repo:       repo,
		nc:         nc,
		userClient: userpb.NewUserServiceClient(userConn),
	}
}

func (uc *CombatUsecase) InitCombat(ctx context.Context, playerID string, playerTeam, enemyTeam []*pb.Character) (string, error) {
	combatID := uuid.New().String()

	// Gather player card IDs to consume energy
	var cardIDs []string
	for _, c := range playerTeam {
		cardIDs = append(cardIDs, c.Id)
	}

	// Call user service to consume energy
	_, err := uc.userClient.ConsumeCardEnergy(ctx, &userpb.ConsumeCardEnergyRequest{
		UserId:  playerID,
		CardIds: cardIDs,
	})
	if err != nil {
		return "", fmt.Errorf("cannot start combat: %v", err)
	}

	// Initialize states (max HP, zero mana)
	var playerState []*pb.CharacterState
	for _, c := range playerTeam {
		playerState = append(playerState, &pb.CharacterState{
			Id:          c.Id,
			CurrentHp:   c.Hp,
			CurrentMana: 0,
		})
	}

	var enemyState []*pb.CharacterState
	for _, c := range enemyTeam {
		enemyState = append(enemyState, &pb.CharacterState{
			Id:          c.Id,
			CurrentHp:   c.Hp,
			CurrentMana: 0,
		})
	}

	state := &domain.CombatState{
		CombatID:        combatID,
		PlayerID:        playerID,
		CurrentTurn:     1,
		PlayerTeam:      playerTeam,
		EnemyTeam:       enemyTeam,
		PlayerTeamState: playerState,
		EnemyTeamState:  enemyState,
		IsEnded:         false,
	}

	err = uc.repo.SaveState(ctx, state)
	if err != nil {
		return "", err
	}

	return combatID, nil
}

func (uc *CombatUsecase) GetCombatState(ctx context.Context, combatID string) (*domain.CombatState, error) {
	return uc.repo.GetState(ctx, combatID)
}

func (uc *CombatUsecase) ExecuteAttack(ctx context.Context, combatID, attackerID, targetID string) (*pb.CombatActionResponse, error) {
	state, err := uc.repo.GetState(ctx, combatID)
	if err != nil {
		return nil, ErrCombatNotFound
	}
	if state.IsEnded {
		return nil, ErrCombatNotFound
	}

	// 1. Find characters
	attacker, _ := findCharacter(attackerID, state)
	target, targetState := findCharacter(targetID, state)

	if attacker == nil || target == nil || targetState == nil {
		return nil, ErrInvalidTarget
	}

	// 2. Calculate Hit/Dodge/Counter
	// Base hit chance
	hitChance := float32(1.0)
	if target.Speed > attacker.Speed {
		hitChance = 1.0 - ((target.Speed - attacker.Speed) / 100.0)
	}
	if hitChance < 0.2 {
		hitChance = 0.2 // Floor of 20%
	}

	isDodged := rand.Float32() > hitChance
	isCountered := false
	damage := int32(0)
	logMsg := ""

	if isDodged {
		// Calculate Counter probability based on Reaction vs Agility
		counterProb := float32(0.0)
		if target.Reaction > attacker.Agility {
			counterProb = float32(target.Reaction-attacker.Agility) / float32(target.Reaction)
			if counterProb > 0.8 {
				counterProb = 0.8 // Max 80% counter chance
			}
		}

		if rand.Float32() < counterProb {
			isCountered = true
			// Target counters attacker
			counterDmg := int32(float32(target.Str) * target.Power) - attacker.Durability
			if counterDmg < 1 {
				counterDmg = 1
			}
			
			// Reduce attacker HP
			_, attackerState := findCharacter(attackerID, state)
			attackerState.CurrentHp -= counterDmg
			if attackerState.CurrentHp < 0 {
				attackerState.CurrentHp = 0
			}
			
			logMsg = fmt.Sprintf("%s dodged and countered %s for %d damage!", target.Name, attacker.Name, counterDmg)
			damage = counterDmg
		} else {
			logMsg = fmt.Sprintf("%s dodged the attack from %s!", target.Name, attacker.Name)
		}
	} else {
		// Normal Attack hits target
		rawDmg := float32(attacker.Str) * attacker.Power
		
		// Critical hit chance (simplified)
		isCritical := rand.Float32() < 0.15 // 15% flat crit chance
		if isCritical {
			rawDmg *= 1.5
		}

		damage = int32(rawDmg) - target.Durability
		if damage < 1 {
			damage = 1
		}

		targetState.CurrentHp -= damage
		if targetState.CurrentHp < 0 {
			targetState.CurrentHp = 0
		}
		
		critText := ""
		if isCritical {
			critText = "CRITICAL HIT! "
		}
		logMsg = fmt.Sprintf("%s%s attacked %s for %d damage.", critText, attacker.Name, target.Name, damage)
	}

	// 3. Check for win condition
	uc.checkWinCondition(state)

	// 4. Save state
	err = uc.repo.SaveState(ctx, state)
	if err != nil {
		return nil, err
	}

	// 5. If ended, publish event
	if state.IsEnded {
		uc.publishCombatEnded(state)
	}

	return &pb.CombatActionResponse{
		Success:      true,
		DamageDealt:  damage,
		ActionLog:    logMsg,
		IsCritical:   false,
		IsDodged:     isDodged,
		IsBlocked:    false,
		IsCountered:  isCountered,
	}, nil
}

// Helpers

func findCharacter(id string, state *domain.CombatState) (*pb.Character, *pb.CharacterState) {
	for i, c := range state.PlayerTeam {
		if c.Id == id {
			return c, state.PlayerTeamState[i]
		}
	}
	for i, c := range state.EnemyTeam {
		if c.Id == id {
			return c, state.EnemyTeamState[i]
		}
	}
	return nil, nil
}

func (uc *CombatUsecase) checkWinCondition(state *domain.CombatState) {
	playerAlive := false
	for _, s := range state.PlayerTeamState {
		if s.CurrentHp > 0 {
			playerAlive = true
			break
		}
	}

	enemyAlive := false
	for _, s := range state.EnemyTeamState {
		if s.CurrentHp > 0 {
			enemyAlive = true
			break
		}
	}

	if !playerAlive {
		state.IsEnded = true
		state.Winner = "enemy"
	} else if !enemyAlive {
		state.IsEnded = true
		state.Winner = "player"
	}
}

func (uc *CombatUsecase) publishCombatEnded(state *domain.CombatState) {
	event := map[string]interface{}{
		"combat_id": state.CombatID,
		"player_id": state.PlayerID,
		"winner":    state.Winner,
	}
	data, _ := json.Marshal(event)
	
	if state.Winner == "player" {
		uc.nc.Publish("combat.won", data)
	} else {
		uc.nc.Publish("combat.lost", data)
	}
}
