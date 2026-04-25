package domain

import (
	"context"
)

// EconomyRepository handles currency and energy persistence (PostgreSQL)
type EconomyRepository interface {
	GetCurrency(ctx context.Context, userID string) (*Currency, error)
	AddCurrency(ctx context.Context, userID string, amount int32) error
	SpendCurrency(ctx context.Context, userID string, amount int32) error

	GetEnergy(ctx context.Context, userID string) (*Energy, error)
	UpdateEnergy(ctx context.Context, energy *Energy) error
}

// GachaCache handles Redis caching for banners and free packs
type GachaCache interface {
	GetBanners(ctx context.Context) ([]*Banner, error)
	GetShopItems(ctx context.Context) ([]*ShopItem, error)
	
	SetFreePackTimer(ctx context.Context, userID string) error
	GetFreePackRemainingSeconds(ctx context.Context, userID string) (int64, error)
}

// QuestRepository handles daily rewards and quests (PostgreSQL)
type QuestRepository interface {
	GetDailyRewards(ctx context.Context, userID string) ([]*DailyReward, error)
	ClaimDailyReward(ctx context.Context, userID string, day int32) error
	
	GetQuests(ctx context.Context, userID string) ([]*Quest, error)
	ClaimQuestReward(ctx context.Context, userID string, questID string) error
}
