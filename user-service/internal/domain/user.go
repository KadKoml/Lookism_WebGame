package domain

import "time"

// User represents a registered player.
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// Card represents a character card owned by a player.
type Card struct {
	ID                   string     `json:"id"`
	UserID               string     `json:"user_id"`
	TemplateID           string     `json:"template_id"`
	Level                int32      `json:"level"`
	MergeStars           int32      `json:"merge_stars"`
	CurrentEnergy        int32      `json:"current_energy"`
	NextRefreshTimestamp *time.Time `json:"next_refresh_timestamp"`
}

// Squad represents the active 3-character team.
type Squad struct {
	UserID  string `json:"user_id"`
	CardID1 string `json:"card_id_1"`
	CardID2 string `json:"card_id_2"`
	CardID3 string `json:"card_id_3"`
}
