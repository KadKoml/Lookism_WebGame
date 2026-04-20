package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/kozie/lookism-rpg/user-service/internal/domain"
)

var (
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound     = errors.New("user not found")
	ErrCardNotFound     = errors.New("card not found")
	ErrMaxLevel         = errors.New("card is already at max level")
)

var jwtSecret = []byte("lookism-rpg-secret-key-2024") // In production, use env var

// UserUsecase handles business logic for users, cards, and squads.
type UserUsecase struct {
	userRepo  domain.UserRepository
	cardRepo  domain.CardRepository
	squadRepo domain.SquadRepository
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(
	userRepo domain.UserRepository,
	cardRepo domain.CardRepository,
	squadRepo domain.SquadRepository,
) *UserUsecase {
	return &UserUsecase{
		userRepo:  userRepo,
		cardRepo:  cardRepo,
		squadRepo: squadRepo,
	}
}

// --- AUTH ---

func (uc *UserUsecase) Register(ctx context.Context, username, email, password string) (*domain.User, string, error) {
	// Check if user already exists
	existing, _ := uc.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, "", ErrUserExists
	}
	existing, _ = uc.userRepo.GetByUsername(ctx, username)
	if existing != nil {
		return nil, "", ErrUserExists
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := uc.generateJWT(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (uc *UserUsecase) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := uc.generateJWT(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (uc *UserUsecase) ValidateToken(tokenStr string) (string, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	userID, ok := claims["user_id"].(string)
	return userID, ok
}

func (uc *UserUsecase) generateJWT(userID, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// --- PROFILE ---

func (uc *UserUsecase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (uc *UserUsecase) UpdateProfile(ctx context.Context, userID, newUsername string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	user.Username = newUsername
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// --- INVENTORY ---

func (uc *UserUsecase) GetInventory(ctx context.Context, userID string) ([]*domain.Card, error) {
	return uc.cardRepo.GetByUserID(ctx, userID)
}

func (uc *UserUsecase) AddCard(ctx context.Context, userID, templateID string) (*domain.Card, error) {
	card := &domain.Card{
		ID:            uuid.New().String(),
		UserID:        userID,
		TemplateID:    templateID,
		Level:         1,
		MergeStars:    0,
		CurrentEnergy: 5,
	}
	if err := uc.cardRepo.Create(ctx, card); err != nil {
		return nil, err
	}
	return card, nil
}

func (uc *UserUsecase) LevelUpCard(ctx context.Context, cardID string) (*domain.Card, error) {
	card, err := uc.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}
	if card.Level >= 60 {
		return nil, ErrMaxLevel
	}
	card.Level++
	if err := uc.cardRepo.Update(ctx, card); err != nil {
		return nil, err
	}
	return card, nil
}

func (uc *UserUsecase) MergeDuplicateCard(ctx context.Context, cardID string) (*domain.Card, error) {
	card, err := uc.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}
	card.MergeStars++
	if err := uc.cardRepo.Update(ctx, card); err != nil {
		return nil, err
	}
	return card, nil
}

func (uc *UserUsecase) GetCardStats(ctx context.Context, cardID string) (*domain.Card, error) {
	card, err := uc.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}
	return card, nil
}

// --- SQUAD ---

func (uc *UserUsecase) SetActiveSquad(ctx context.Context, userID string, cardIDs []string) error {
	squad := &domain.Squad{UserID: userID}
	if len(cardIDs) > 0 {
		squad.CardID1 = cardIDs[0]
	}
	if len(cardIDs) > 1 {
		squad.CardID2 = cardIDs[1]
	}
	if len(cardIDs) > 2 {
		squad.CardID3 = cardIDs[2]
	}
	return uc.squadRepo.Upsert(ctx, squad)
}

func (uc *UserUsecase) GetActiveSquad(ctx context.Context, userID string) (*domain.Squad, []*domain.Card, error) {
	squad, err := uc.squadRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	var cards []*domain.Card
	for _, cid := range []string{squad.CardID1, squad.CardID2, squad.CardID3} {
		if cid != "" {
			card, err := uc.cardRepo.GetByID(ctx, cid)
			if err == nil {
				cards = append(cards, card)
			}
		}
	}
	return squad, cards, nil
}

func (uc *UserUsecase) ConsumeCardEnergy(ctx context.Context, userID string, cardIDs []string) error {
	for _, cid := range cardIDs {
		card, err := uc.cardRepo.GetByID(ctx, cid)
		if err != nil {
			return err
		}
		if card.UserID != userID {
			return errors.New("card does not belong to user")
		}
		if card.CurrentEnergy <= 0 {
			return fmt.Errorf("card %s has no energy left", card.ID)
		}
		
		card.CurrentEnergy--
		
		// If it just dropped from full, start the refresh timer
		if card.CurrentEnergy == 4 || card.NextRefreshTimestamp == nil {
			now := time.Now().Add(5 * time.Minute) // 5 minutes to refresh 1 energy
			card.NextRefreshTimestamp = &now
		}

		if err := uc.cardRepo.Update(ctx, card); err != nil {
			return err
		}
	}
	return nil
}
