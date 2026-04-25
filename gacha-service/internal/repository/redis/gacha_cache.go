package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/kozie/lookism-rpg/gacha-service/internal/domain"
)

type gachaCache struct {
	client *redis.Client
}

func NewGachaCache(client *redis.Client) domain.GachaCache {
	return &gachaCache{client: client}
}

func (c *gachaCache) GetBanners(ctx context.Context) ([]*domain.Banner, error) {
	val, err := c.client.Get(ctx, "cache:banners").Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	var banners []*domain.Banner
	err = json.Unmarshal([]byte(val), &banners)
	return banners, err
}

func (c *gachaCache) GetShopItems(ctx context.Context) ([]*domain.ShopItem, error) {
	val, err := c.client.Get(ctx, "cache:shop_items").Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	var items []*domain.ShopItem
	err = json.Unmarshal([]byte(val), &items)
	return items, err
}

func (c *gachaCache) SetFreePackTimer(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:%s:free_pack_timer", userID)
	// 1 minute cooldown for free pack testing
	return c.client.Set(ctx, key, "locked", time.Minute).Err()
}

func (c *gachaCache) GetFreePackRemainingSeconds(ctx context.Context, userID string) (int64, error) {
	key := fmt.Sprintf("user:%s:free_pack_timer", userID)
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	
	// If TTL is negative, the key expired (ready) or doesn't exist (ready)
	if ttl <= 0 {
		return 0, nil
	}
	
	return int64(ttl.Seconds()), nil
}

// Helper methods to populate cache (normally called by some admin or background job)
func (c *gachaCache) SetBanners(ctx context.Context, banners []*domain.Banner) error {
	data, err := json.Marshal(banners)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "cache:banners", data, 24*time.Hour).Err()
}

func (c *gachaCache) SetShopItems(ctx context.Context, items []*domain.ShopItem) error {
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "cache:shop_items", data, 24*time.Hour).Err()
}
