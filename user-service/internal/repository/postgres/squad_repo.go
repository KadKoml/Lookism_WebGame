package postgres

import (
	"context"
	"database/sql"

	"github.com/kozie/lookism-rpg/user-service/internal/domain"
)

type squadRepo struct {
	db *sql.DB
}

// NewSquadRepository creates a new PostgreSQL-backed squad repository.
func NewSquadRepository(db *sql.DB) domain.SquadRepository {
	return &squadRepo{db: db}
}

func (r *squadRepo) Upsert(ctx context.Context, squad *domain.Squad) error {
	query := `INSERT INTO squads (user_id, card_id_1, card_id_2, card_id_3)
	           VALUES ($1, $2, $3, $4)
	           ON CONFLICT (user_id) DO UPDATE
	           SET card_id_1 = EXCLUDED.card_id_1,
	               card_id_2 = EXCLUDED.card_id_2,
	               card_id_3 = EXCLUDED.card_id_3`
	_, err := r.db.ExecContext(ctx, query,
		squad.UserID, squad.CardID1, squad.CardID2, squad.CardID3)
	return err
}

func (r *squadRepo) GetByUserID(ctx context.Context, userID string) (*domain.Squad, error) {
	squad := &domain.Squad{}
	query := `SELECT user_id, card_id_1, card_id_2, card_id_3 FROM squads WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&squad.UserID, &squad.CardID1, &squad.CardID2, &squad.CardID3)
	if err != nil {
		return nil, err
	}
	return squad, nil
}
