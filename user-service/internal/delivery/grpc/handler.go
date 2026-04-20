package grpc

import (
	"context"

	pb "github.com/kozie/lookism-rpg/api/proto/user"
	"github.com/kozie/lookism-rpg/user-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler implements the gRPC UserServiceServer.
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	uc *usecase.UserUsecase
}

// NewUserHandler creates a new gRPC handler.
func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

// --- AUTH ---

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	user, token, err := h.uc.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		if err == usecase.ErrUserExists {
			return &pb.AuthResponse{Success: false}, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Errorf(codes.Internal, "registration failed: %v", err)
	}
	return &pb.AuthResponse{
		UserId:   user.ID,
		Username: user.Username,
		Token:    token,
		Success:  true,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user, token, err := h.uc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &pb.AuthResponse{Success: false}, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}
	return &pb.AuthResponse{
		UserId:   user.ID,
		Username: user.Username,
		Token:    token,
		Success:  true,
	}, nil
}

func (h *UserHandler) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	userID, valid := h.uc.ValidateToken(req.Token)
	return &pb.TokenResponse{
		IsValid: valid,
		UserId:  userID,
	}, nil
}

// --- PROFILE ---

func (h *UserHandler) GetProfile(ctx context.Context, req *pb.UserRequest) (*pb.ProfileResponse, error) {
	user, err := h.uc.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &pb.ProfileResponse{
		UserId:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Unix(),
		Success:   true,
	}, nil
}

func (h *UserHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.ProfileResponse, error) {
	user, err := h.uc.UpdateProfile(ctx, req.UserId, req.NewUsername)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update failed: %v", err)
	}
	return &pb.ProfileResponse{
		UserId:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Unix(),
		Success:   true,
	}, nil
}

// --- INVENTORY ---

func (h *UserHandler) GetInventory(ctx context.Context, req *pb.UserRequest) (*pb.InventoryResponse, error) {
	cards, err := h.uc.GetInventory(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get inventory: %v", err)
	}
	var pbCards []*pb.Card
	for _, c := range cards {
		pbCards = append(pbCards, &pb.Card{
			Id:         c.ID,
			TemplateId: c.TemplateID,
			Level:      c.Level,
			MergeStars: c.MergeStars,
		})
	}
	return &pb.InventoryResponse{Cards: pbCards}, nil
}

func (h *UserHandler) AddCardToInventory(ctx context.Context, req *pb.AddCardRequest) (*pb.CardResponse, error) {
	card, err := h.uc.AddCard(ctx, req.UserId, req.TemplateId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add card: %v", err)
	}
	return &pb.CardResponse{
		Card: &pb.Card{
			Id:         card.ID,
			TemplateId: card.TemplateID,
			Level:      card.Level,
			MergeStars: card.MergeStars,
		},
		Success: true,
		Message: "Card added to inventory",
	}, nil
}

// --- CARD PROGRESSION ---

func (h *UserHandler) LevelUpCard(ctx context.Context, req *pb.LevelUpRequest) (*pb.CardResponse, error) {
	card, err := h.uc.LevelUpCard(ctx, req.CardId)
	if err != nil {
		if err == usecase.ErrMaxLevel {
			return &pb.CardResponse{Success: false, Message: "Card is at max level"}, nil
		}
		return nil, status.Errorf(codes.Internal, "level up failed: %v", err)
	}
	return &pb.CardResponse{
		Card: &pb.Card{
			Id:         card.ID,
			TemplateId: card.TemplateID,
			Level:      card.Level,
			MergeStars: card.MergeStars,
		},
		Success: true,
		Message: "Card leveled up!",
	}, nil
}

func (h *UserHandler) MergeDuplicateCard(ctx context.Context, req *pb.MergeRequest) (*pb.CardResponse, error) {
	card, err := h.uc.MergeDuplicateCard(ctx, req.CardId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "merge failed: %v", err)
	}
	return &pb.CardResponse{
		Card: &pb.Card{
			Id:         card.ID,
			TemplateId: card.TemplateID,
			Level:      card.Level,
			MergeStars: card.MergeStars,
		},
		Success: true,
		Message: "Card merged! +1 star",
	}, nil
}

func (h *UserHandler) GetCardStats(ctx context.Context, req *pb.CardRequest) (*pb.CardStatsResponse, error) {
	card, err := h.uc.GetCardStats(ctx, req.CardId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "card not found")
	}

	// Use the existing CalculateFinalIntStat helper
	finalHP := usecase.CalculateFinalIntStat(430, card.Level, card.MergeStars)   // Example base
	finalSTR := usecase.CalculateFinalIntStat(210, card.Level, card.MergeStars)  // Example base

	return &pb.CardStatsResponse{
		CardId:          card.ID,
		Name:            card.TemplateID,
		FinalHp:         finalHP,
		FinalStr:        finalSTR,
		FinalMana:       4,
		FinalAgility:    usecase.CalculateFinalIntStat(48, card.Level, card.MergeStars),
		FinalReaction:   usecase.CalculateFinalIntStat(50, card.Level, card.MergeStars),
		FinalDurability: usecase.CalculateFinalIntStat(42, card.Level, card.MergeStars),
		FinalPower:      float32(usecase.CalculateFinalStat(1.35, card.Level, card.MergeStars)),
		FinalSpeed:      float32(usecase.CalculateFinalStat(1.35, card.Level, card.MergeStars)),
		FinalTechnique:  float32(usecase.CalculateFinalStat(1.45, card.Level, card.MergeStars)),
	}, nil
}

// --- SQUAD ---

func (h *UserHandler) SetActiveSquad(ctx context.Context, req *pb.SquadRequest) (*pb.SquadResponse, error) {
	err := h.uc.SetActiveSquad(ctx, req.UserId, req.CardIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set squad: %v", err)
	}
	return &pb.SquadResponse{
		UserId:  req.UserId,
		Success: true,
	}, nil
}

func (h *UserHandler) GetActiveSquad(ctx context.Context, req *pb.UserRequest) (*pb.SquadResponse, error) {
	_, cards, err := h.uc.GetActiveSquad(ctx, req.UserId)
	if err != nil {
		return &pb.SquadResponse{UserId: req.UserId, Success: false}, nil
	}
	var pbCards []*pb.Card
	for _, c := range cards {
		pbCards = append(pbCards, &pb.Card{
			Id:         c.ID,
			TemplateId: c.TemplateID,
			Level:      c.Level,
			MergeStars: c.MergeStars,
		})
	}
	return &pb.SquadResponse{
		UserId:  req.UserId,
		Squad:   pbCards,
		Success: true,
	}, nil
}

func (h *UserHandler) ConsumeCardEnergy(ctx context.Context, req *pb.ConsumeCardEnergyRequest) (*pb.ConsumeCardEnergyResponse, error) {
	err := h.uc.ConsumeCardEnergy(ctx, req.UserId, req.CardIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to consume energy: %v", err)
	}
	return &pb.ConsumeCardEnergyResponse{
		Success: true,
		Message: "Energy consumed successfully",
	}, nil
}
