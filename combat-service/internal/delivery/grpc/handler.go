package grpc

import (
	"context"

	pb "github.com/kozie/lookism-rpg/api/proto/combat"
	"github.com/kozie/lookism-rpg/combat-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CombatHandler struct {
	pb.UnimplementedCombatEngineServiceServer
	uc *usecase.CombatUsecase
}

func NewCombatHandler(uc *usecase.CombatUsecase) *CombatHandler {
	return &CombatHandler{uc: uc}
}

func (h *CombatHandler) InitCombatState(ctx context.Context, req *pb.InitCombatRequest) (*pb.InitCombatResponse, error) {
	combatID, err := h.uc.InitCombat(ctx, req.PlayerId, req.PlayerTeam, req.EnemyTeam)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to init combat: %v", err)
	}

	return &pb.InitCombatResponse{
		CombatId: combatID,
		Success:  true,
	}, nil
}

func (h *CombatHandler) GetCombatState(ctx context.Context, req *pb.CombatRequest) (*pb.CombatStateResponse, error) {
	state, err := h.uc.GetCombatState(ctx, req.CombatId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "combat state not found: %v", err)
	}

	return &pb.CombatStateResponse{
		CombatId:        state.CombatID,
		CurrentTurn:     state.CurrentTurn,
		PlayerTeamState: state.PlayerTeamState,
		EnemyTeamState:  state.EnemyTeamState,
	}, nil
}

func (h *CombatHandler) ExecuteAttack(ctx context.Context, req *pb.CombatActionRequest) (*pb.CombatActionResponse, error) {
	resp, err := h.uc.ExecuteAttack(ctx, req.CombatId, req.AttackerId, req.TargetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "attack failed: %v", err)
	}
	return resp, nil
}

// Dummy implementations for other proto methods to satisfy interface
func (h *CombatHandler) CampaignMatchmaking(ctx context.Context, req *pb.MatchmakingRequest) (*pb.MatchmakingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CampaignMatchmaking not implemented")
}

func (h *CombatHandler) GetCombatFormulas(ctx context.Context, req *pb.Empty) (*pb.FormulasResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCombatFormulas not implemented")
}

func (h *CombatHandler) ExecuteBlock(ctx context.Context, req *pb.CombatActionRequest) (*pb.CombatActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteBlock not implemented")
}

func (h *CombatHandler) ExecuteDodge(ctx context.Context, req *pb.CombatActionRequest) (*pb.CombatActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteDodge not implemented")
}

func (h *CombatHandler) ExecuteCounter(ctx context.Context, req *pb.CombatActionRequest) (*pb.CombatActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteCounter not implemented")
}

func (h *CombatHandler) ActivateActiveSkill(ctx context.Context, req *pb.SkillActionRequest) (*pb.CombatActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivateActiveSkill not implemented")
}

func (h *CombatHandler) ApplyTeamPassives(ctx context.Context, req *pb.CombatRequest) (*pb.CombatStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApplyTeamPassives not implemented")
}

func (h *CombatHandler) CalculateTurnOutcome(ctx context.Context, req *pb.CombatRequest) (*pb.TurnOutcomeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculateTurnOutcome not implemented")
}

func (h *CombatHandler) ApplyTickEffects(ctx context.Context, req *pb.CombatRequest) (*pb.CombatStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApplyTickEffects not implemented")
}

func (h *CombatHandler) ForceEndCombat(ctx context.Context, req *pb.CombatRequest) (*pb.EndCombatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForceEndCombat not implemented")
}
