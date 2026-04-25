package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/kozie/lookism-rpg/gacha-service/internal/domain"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type economyRepo struct {
	db *sql.DB
}

func NewEconomyRepository(db *sql.DB) domain.EconomyRepository {
	return &economyRepo{db: db}
}

func (r *economyRepo) GetCurrency(ctx context.Context, userID string) (*domain.Currency, error) {
	currency := &domain.Currency{UserID: userID}
	query := `SELECT balance FROM currency WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&currency.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			// Auto-initialize with 1000 starting currency if not found
			_, errInit := r.db.ExecContext(ctx, `INSERT INTO currency (user_id, balance) VALUES ($1, $2)`, userID, 1000)
			if errInit != nil {
				return nil, errInit
			}
			currency.Balance = 1000
			return currency, nil
		}
		return nil, err
	}
	return currency, nil
}

func (r *economyRepo) AddCurrency(ctx context.Context, userID string, amount int32) error {
	query := `INSERT INTO currency (user_id, balance) VALUES ($1, $2)
	          ON CONFLICT (user_id) DO UPDATE SET balance = currency.balance + EXCLUDED.balance`
	_, err := r.db.ExecContext(ctx, query, userID, amount)
	return err
}

func (r *economyRepo) SpendCurrency(ctx context.Context, userID string, amount int32) error {
	// Need to check balance first to avoid going negative
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance int32
	err = tx.QueryRowContext(ctx, `SELECT balance FROM currency WHERE user_id = $1 FOR UPDATE`, userID).Scan(&currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrInsufficientFunds
		}
		return err
	}

	if currentBalance < amount {
		return ErrInsufficientFunds
	}

	_, err = tx.ExecContext(ctx, `UPDATE currency SET balance = balance - $1 WHERE user_id = $2`, amount, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *economyRepo) GetEnergy(ctx context.Context, userID string) (*domain.Energy, error) {
	energy := &domain.Energy{UserID: userID}
	query := `SELECT current_energy, max_energy, last_refill FROM energy WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&energy.CurrentEnergy, &energy.MaxEnergy, &energy.LastRefill)
	if err != nil {
		if err == sql.ErrNoRows {
			// Auto-init
			energy.CurrentEnergy = 40
			energy.MaxEnergy = 40
			energy.LastRefill = time.Now()
			_, errInit := r.db.ExecContext(ctx, `INSERT INTO energy (user_id, current_energy, max_energy, last_refill) VALUES ($1, $2, $3, $4)`,
				userID, energy.CurrentEnergy, energy.MaxEnergy, energy.LastRefill)
			if errInit != nil {
				return nil, errInit
			}
			return energy, nil
		}
		return nil, err
	}
	return energy, nil
}

func (r *economyRepo) UpdateEnergy(ctx context.Context, energy *domain.Energy) error {
	query := `UPDATE energy SET current_energy = $1, last_refill = $2 WHERE user_id = $3`
	_, err := r.db.ExecContext(ctx, query, energy.CurrentEnergy, energy.LastRefill, energy.UserID)
	return err
}
