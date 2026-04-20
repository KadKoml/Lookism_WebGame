package postgres

import (
	"context"
	"database/sql"

	"github.com/kozie/lookism-rpg/user-service/internal/domain"
)

type cardRepo struct {
	db *sql.DB
}

// NewCardRepository creates a new PostgreSQL-backed card repository.
func NewCardRepository(db *sql.DB) domain.CardRepository {
	return &cardRepo{db: db}
}

func (r *cardRepo) Create(ctx context.Context, card *domain.Card) error {
	query := `INSERT INTO cards (id, user_id, template_id, level, merge_stars, current_energy)
	           VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		card.ID, card.UserID, card.TemplateID, card.Level, card.MergeStars, card.CurrentEnergy)
	return err
}

func (r *cardRepo) GetByID(ctx context.Context, id string) (*domain.Card, error) {
	card := &domain.Card{}
	query := `SELECT id, user_id, template_id, level, merge_stars, current_energy, next_refresh_timestamp FROM cards WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&card.ID, &card.UserID, &card.TemplateID, &card.Level, &card.MergeStars, &card.CurrentEnergy, &card.NextRefreshTimestamp)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (r *cardRepo) GetByUserID(ctx context.Context, userID string) ([]*domain.Card, error) {
	query := `SELECT id, user_id, template_id, level, merge_stars, current_energy, next_refresh_timestamp FROM cards WHERE user_id = $1 ORDER BY level DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*domain.Card
	for rows.Next() {
		card := &domain.Card{}
		if err := rows.Scan(&card.ID, &card.UserID, &card.TemplateID, &card.Level, &card.MergeStars, &card.CurrentEnergy, &card.NextRefreshTimestamp); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func (r *cardRepo) Update(ctx context.Context, card *domain.Card) error {
	query := `UPDATE cards SET level = $1, merge_stars = $2, current_energy = $3, next_refresh_timestamp = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, card.Level, card.MergeStars, card.CurrentEnergy, card.NextRefreshTimestamp, card.ID)
	return err
}

func (r *cardRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cards WHERE id = $1`, id)
	return err
}
