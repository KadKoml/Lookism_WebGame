package grpc

import (
	"context"

	pb "github.com/kozie/lookism-rpg/api/proto/gacha"
	"github.com/kozie/lookism-rpg/gacha-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GachaHandler struct {
	pb.UnimplementedGachaEconomyServiceServer
	uc *usecase.GachaUsecase
}

func NewGachaHandler(uc *usecase.GachaUsecase) *GachaHandler {
	return &GachaHandler{uc: uc}
}

func (h *GachaHandler) GetBanners(ctx context.Context, req *pb.Empty) (*pb.BannerList, error) {
	banners, err := h.uc.GetBanners(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get banners: %v", err)
	}

	var pbBanners []*pb.Banner
	for _, b := range banners {
		pbBanners = append(pbBanners, &pb.Banner{
			Id:   b.ID,
			Name: b.Name,
			Cost: b.Cost,
		})
	}
	return &pb.BannerList{Banners: pbBanners}, nil
}

func (h *GachaHandler) RollGachaSingle(ctx context.Context, req *pb.RollRequest) (*pb.RollResponse, error) {
	cards, err := h.uc.RollGacha(ctx, req.UserId, req.BannerId, 1)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "roll failed: %v", err)
	}

	var pulled []*pb.PulledCard
	for _, c := range cards {
		pulled = append(pulled, &pb.PulledCard{
			CardId:        c.Id,
			TemplateId:    c.TemplateId,
			NewMergeStars: c.MergeStars,
		})
	}

	return &pb.RollResponse{
		Cards:   pulled,
		Success: true,
		Message: "Rolled successfully!",
	}, nil
}

func (h *GachaHandler) RollGachaMulti(ctx context.Context, req *pb.RollRequest) (*pb.RollResponse, error) {
	cards, err := h.uc.RollGacha(ctx, req.UserId, req.BannerId, 10)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "roll failed: %v", err)
	}

	var pulled []*pb.PulledCard
	for _, c := range cards {
		pulled = append(pulled, &pb.PulledCard{
			CardId:        c.Id,
			TemplateId:    c.TemplateId,
			NewMergeStars: c.MergeStars,
		})
	}

	return &pb.RollResponse{
		Cards:   pulled,
		Success: true,
		Message: "10x Roll successful!",
	}, nil
}

func (h *GachaHandler) AddCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.CurrencyResponse, error) {
	err := h.uc.AddCurrency(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add currency: %v", err)
	}
	
	curr, _ := h.uc.GetCurrency(ctx, req.UserId)
	balance := int32(0)
	if curr != nil {
		balance = curr.Balance
	}

	return &pb.CurrencyResponse{
		Success:    true,
		NewBalance: balance,
		Message:    "Currency added",
	}, nil
}

func (h *GachaHandler) SpendCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.CurrencyResponse, error) {
	err := h.uc.SpendCurrency(ctx, req.UserId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to spend currency: %v", err)
	}
	
	curr, _ := h.uc.GetCurrency(ctx, req.UserId)
	balance := int32(0)
	if curr != nil {
		balance = curr.Balance
	}

	return &pb.CurrencyResponse{
		Success:    true,
		NewBalance: balance,
		Message:    "Currency spent",
	}, nil
}

func (h *GachaHandler) GetDailyRewards(ctx context.Context, req *pb.UserRequest) (*pb.DailyRewardsList, error) {
	rewards, err := h.uc.GetDailyRewards(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get rewards: %v", err)
	}

	var pbRewards []*pb.DailyReward
	for _, r := range rewards {
		pbRewards = append(pbRewards, &pb.DailyReward{
			Day:         r.Day,
			Description: r.Description,
			Claimed:     r.Claimed,
		})
	}
	return &pb.DailyRewardsList{Rewards: pbRewards}, nil
}

func (h *GachaHandler) ClaimDailyReward(ctx context.Context, req *pb.ClaimRequest) (*pb.ClaimResponse, error) {
	// Parse req.TargetId as Day int. For simplicity, just assuming day 1
	err := h.uc.ClaimDailyReward(ctx, req.UserId, 1) 
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to claim: %v", err)
	}
	return &pb.ClaimResponse{Success: true, Message: "Reward claimed"}, nil
}

func (h *GachaHandler) GetQuests(ctx context.Context, req *pb.UserRequest) (*pb.QuestList, error) {
	quests, err := h.uc.GetQuests(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get quests: %v", err)
	}

	var pbQuests []*pb.Quest
	for _, q := range quests {
		pbQuests = append(pbQuests, &pb.Quest{
			Id:          q.ID,
			Description: q.Description,
			IsCompleted: q.IsCompleted,
			IsClaimed:   q.IsClaimed,
		})
	}
	return &pb.QuestList{Quests: pbQuests}, nil
}

func (h *GachaHandler) ClaimQuestReward(ctx context.Context, req *pb.ClaimRequest) (*pb.ClaimResponse, error) {
	err := h.uc.ClaimQuestReward(ctx, req.UserId, req.TargetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to claim quest: %v", err)
	}
	return &pb.ClaimResponse{Success: true, Message: "Quest claimed"}, nil
}

func (h *GachaHandler) GetFreePackTimer(ctx context.Context, req *pb.UserRequest) (*pb.TimerResponse, error) {
	remaining, err := h.uc.GetFreePackTimer(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get timer: %v", err)
	}
	return &pb.TimerResponse{
		RemainingSeconds: remaining,
		IsReady:          remaining <= 0,
	}, nil
}
