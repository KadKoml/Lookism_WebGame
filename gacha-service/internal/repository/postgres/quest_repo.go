package postgres

import (
	"context"
	"database/sql"

	"github.com/kozie/lookism-rpg/gacha-service/internal/domain"
)

type questRepo struct {
	db *sql.DB
}

func NewQuestRepository(db *sql.DB) domain.QuestRepository {
	return &questRepo{db: db}
}

func (r *questRepo) GetDailyRewards(ctx context.Context, userID string) ([]*domain.DailyReward, error) {
	query := `SELECT day, description, claimed FROM daily_rewards WHERE user_id = $1 ORDER BY day ASC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rewards []*domain.DailyReward
	for rows.Next() {
		reward := &domain.DailyReward{UserID: userID}
		if err := rows.Scan(&reward.Day, &reward.Description, &reward.Claimed); err != nil {
			return nil, err
		}
		rewards = append(rewards, reward)
	}

	// Initialize default 7-day rewards if empty
	if len(rewards) == 0 {
		return r.initDailyRewards(ctx, userID)
	}

	return rewards, nil
}

func (r *questRepo) initDailyRewards(ctx context.Context, userID string) ([]*domain.DailyReward, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	defaultRewards := []struct {
		day  int32
		desc string
	}{
		{1, "100 Монет"},
		{2, "200 Монет"},
		{3, "300 Монет"},
		{4, "400 Монет"},
		{5, "500 Монет"},
		{6, "600 Монет"},
		{7, "Случайный персонаж"},
	}

	var rewards []*domain.DailyReward
	for _, dr := range defaultRewards {
		_, err := tx.ExecContext(ctx, `INSERT INTO daily_rewards (user_id, day, description, claimed) VALUES ($1, $2, $3, false)`, userID, dr.day, dr.desc)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, &domain.DailyReward{UserID: userID, Day: dr.day, Description: dr.desc, Claimed: false})
	}

	return rewards, tx.Commit()
}

func (r *questRepo) ClaimDailyReward(ctx context.Context, userID string, day int32) error {
	query := `UPDATE daily_rewards SET claimed = true WHERE user_id = $1 AND day = $2`
	_, err := r.db.ExecContext(ctx, query, userID, day)
	return err
}

func (r *questRepo) GetQuests(ctx context.Context, userID string) ([]*domain.Quest, error) {
	query := `SELECT id, description, is_completed, is_claimed FROM quests WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quests []*domain.Quest
	for rows.Next() {
		q := &domain.Quest{UserID: userID}
		if err := rows.Scan(&q.ID, &q.Description, &q.IsCompleted, &q.IsClaimed); err != nil {
			return nil, err
		}
		quests = append(quests, q)
	}
	return quests, nil
}

func (r *questRepo) ClaimQuestReward(ctx context.Context, userID string, questID string) error {
	query := `UPDATE quests SET is_claimed = true WHERE user_id = $1 AND id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, questID)
	return err
}
