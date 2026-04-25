package domain

import "time"

type Banner struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Cost int32  `json:"cost"`
}

type Currency struct {
	UserID  string `json:"user_id"`
	Balance int32  `json:"balance"`
}

type Energy struct {
	UserID        string    `json:"user_id"`
	CurrentEnergy int32     `json:"current_energy"`
	MaxEnergy     int32     `json:"max_energy"`
	LastRefill    time.Time `json:"last_refill"`
}

type DailyReward struct {
	UserID      string `json:"user_id"`
	Day         int32  `json:"day"`
	Description string `json:"description"`
	Claimed     bool   `json:"claimed"`
}

type Quest struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	IsClaimed   bool   `json:"is_claimed"`
}

type ShopItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Cost        int32  `json:"cost"`
	Description string `json:"description"`
}
