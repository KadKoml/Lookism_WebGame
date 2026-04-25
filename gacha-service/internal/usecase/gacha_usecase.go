package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	"github.com/kozie/lookism-rpg/gacha-service/internal/domain"
	userpb "github.com/kozie/lookism-rpg/api/proto/user"
)

type FreePackEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	ErrBannerNotFound = errors.New("banner not found")
	ErrFreePackNotReady = errors.New("free pack is not ready yet")
)

type GachaUsecase struct {
	economyRepo domain.EconomyRepository
	questRepo   domain.QuestRepository
	cache       domain.GachaCache
	nc          *nats.Conn
	userClient  userpb.UserServiceClient
}

func NewGachaUsecase(
	economyRepo domain.EconomyRepository,
	questRepo domain.QuestRepository,
	cache domain.GachaCache,
	nc *nats.Conn,
	userConn *grpc.ClientConn,
) *GachaUsecase {
	return &GachaUsecase{
		economyRepo: economyRepo,
		questRepo:   questRepo,
		cache:       cache,
		nc:          nc,
		userClient:  userpb.NewUserServiceClient(userConn),
	}
}

// Default Banners fallback
var defaultBanners = []*domain.Banner{
	{ID: "basic_banner", Name: "Basic Lookism Pack", Cost: 500},
	{ID: "premium_banner", Name: "Premium Lookism Pack", Cost: 1000},
}

// Available card templates (Hardcoded for simplicity, usually from DB)
var availableTemplates = []string{"c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8", "c9", "c10", "c11", "c12", "c13", "c14", "c15"}

// --- GACHA & SHOP ---

func (uc *GachaUsecase) GetBanners(ctx context.Context) ([]*domain.Banner, error) {
	banners, err := uc.cache.GetBanners(ctx)
	if err != nil || len(banners) == 0 {
		return defaultBanners, nil
	}
	return banners, nil
}

func (uc *GachaUsecase) RollGacha(ctx context.Context, userID, bannerID string, count int) ([]*userpb.Card, error) {
	// 1. Find banner cost
	banners, _ := uc.GetBanners(ctx)
	var banner *domain.Banner
	for _, b := range banners {
		if b.ID == bannerID {
			banner = b
			break
		}
	}
	if banner == nil {
		return nil, ErrBannerNotFound
	}

	totalCost := banner.Cost * int32(count)

	// 2. Spend currency
	err := uc.economyRepo.SpendCurrency(ctx, userID, totalCost)
	if err != nil {
		return nil, err
	}

	// 3. Roll cards and call user-service
	var pulledCards []*userpb.Card
	for i := 0; i < count; i++ {
		templateID := availableTemplates[rand.Intn(len(availableTemplates))]
		
		// Call user-service to add card
		resp, err := uc.userClient.AddCardToInventory(ctx, &userpb.AddCardRequest{
			UserId:     userID,
			TemplateId: templateID,
		})
		if err != nil {
			// In a real saga, we would compensate (refund currency) here if adding card fails
			return nil, err
		}
		pulledCards = append(pulledCards, resp.Card)
	}

	// 4. Publish NATS event
	eventData, _ := json.Marshal(map[string]interface{}{
		"user_id": userID,
		"count":   count,
	})
	uc.nc.Publish("gacha.roll_completed", eventData)

	return pulledCards, nil
}

func (uc *GachaUsecase) GetFreePackTimer(ctx context.Context, userID string) (int64, error) {
	return uc.cache.GetFreePackRemainingSeconds(ctx, userID)
}

func (uc *GachaUsecase) ClaimFreePack(ctx context.Context, userID string) (*userpb.Card, error) {
	remaining, err := uc.GetFreePackTimer(ctx, userID)
	if err != nil {
		return nil, err
	}
	if remaining > 0 {
		return nil, ErrFreePackNotReady
	}

	// Reset timer
	uc.cache.SetFreePackTimer(ctx, userID)

	// Roll 1 card
	templateID := availableTemplates[rand.Intn(len(availableTemplates))]
	resp, err := uc.userClient.AddCardToInventory(ctx, &userpb.AddCardRequest{
		UserId:     userID,
		TemplateId: templateID,
	})
	if err != nil {
		return nil, err
	}
	return resp.Card, nil
}

// CheckAndPublishFreePack checks if 1 hour of playtime has passed since the last claim.
// If so, it publishes a FreePackReady event to NATS.
func (uc *GachaUsecase) CheckAndPublishFreePack(ctx context.Context, userID, userEmail string, lastClaimTime time.Time) error {
	// Simple check: 1 hour has passed
	if time.Since(lastClaimTime) >= time.Hour {
		event := FreePackEvent{
			UserID:    userID,
			Email:     userEmail,
			Timestamp: time.Now(),
		}
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}
		
		// Publish event to Notification service
		err = uc.nc.Publish("notifications.free_pack", data)
		if err != nil {
			return err
		}
	}
	return nil
}

// --- ECONOMY ---

func (uc *GachaUsecase) GetCurrency(ctx context.Context, userID string) (*domain.Currency, error) {
	return uc.economyRepo.GetCurrency(ctx, userID)
}

func (uc *GachaUsecase) AddCurrency(ctx context.Context, userID string, amount int32) error {
	return uc.economyRepo.AddCurrency(ctx, userID, amount)
}

func (uc *GachaUsecase) SpendCurrency(ctx context.Context, userID string, amount int32) error {
	return uc.economyRepo.SpendCurrency(ctx, userID, amount)
}

// --- QUESTS & REWARDS ---

func (uc *GachaUsecase) GetDailyRewards(ctx context.Context, userID string) ([]*domain.DailyReward, error) {
	return uc.questRepo.GetDailyRewards(ctx, userID)
}

func (uc *GachaUsecase) ClaimDailyReward(ctx context.Context, userID string, day int32) error {
	// Mark claimed
	err := uc.questRepo.ClaimDailyReward(ctx, userID, day)
	if err != nil {
		return err
	}
	// Add currency (assuming reward is currency for simplicity)
	return uc.economyRepo.AddCurrency(ctx, userID, 100 * day) 
}

func (uc *GachaUsecase) GetQuests(ctx context.Context, userID string) ([]*domain.Quest, error) {
	return uc.questRepo.GetQuests(ctx, userID)
}

func (uc *GachaUsecase) ClaimQuestReward(ctx context.Context, userID string, questID string) error {
	err := uc.questRepo.ClaimQuestReward(ctx, userID, questID)
	if err != nil {
		return err
	}
	return uc.economyRepo.AddCurrency(ctx, userID, 500) // Default quest reward
}
