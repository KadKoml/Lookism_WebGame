package consumer

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/kozie/lookism-rpg/notification-service/internal/service"
)

type EnergyWorker struct {
	db     *sql.DB
	sender service.SMTPSender
}

func NewEnergyWorker(db *sql.DB, sender service.SMTPSender) *EnergyWorker {
	return &EnergyWorker{
		db:     db,
		sender: sender,
	}
}

func (w *EnergyWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	log.Println("Started Energy Worker to monitor card refreshes...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Energy Worker...")
			return
		case <-ticker.C:
			w.checkRefreshedCards(ctx)
		}
	}
}

func (w *EnergyWorker) checkRefreshedCards(ctx context.Context) {
	// Find cards where next_refresh_timestamp has passed AND energy is not full
	// Wait, we need to know the user's email to notify them. 
	// We have to join the cards table with the users table.
	query := `
		SELECT c.id, c.template_id, c.current_energy, u.id, u.email
		FROM cards c
		JOIN users u ON c.user_id = u.id
		WHERE c.next_refresh_timestamp IS NOT NULL
		  AND c.next_refresh_timestamp <= NOW()
		  AND c.current_energy < 5
	`
	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("EnergyWorker DB error: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cardID, templateID, userID, userEmail string
		var currentEnergy int
		if err := rows.Scan(&cardID, &templateID, &currentEnergy, &userID, &userEmail); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		// Update card energy
		newEnergy := currentEnergy + 1
		var nextRefresh *time.Time
		if newEnergy < 5 {
			t := time.Now().Add(5 * time.Minute)
			nextRefresh = &t
		}

		// Execute update
		_, err = w.db.ExecContext(ctx, "UPDATE cards SET current_energy = $1, next_refresh_timestamp = $2 WHERE id = $3", newEnergy, nextRefresh, cardID)
		if err != nil {
			log.Printf("Failed to update card energy: %v", err)
			continue
		}

		log.Printf("Card %s recovered 1 energy! Now at %d/5", cardID, newEnergy)

		// If fully refreshed, send email
		if newEnergy == 5 {
			log.Printf("Card %s is fully refreshed! Sending email to %s", cardID, userEmail)
			
			subject := "Lookism RPG: Your Card is Ready for Battle!"
			body := "Your card is fully rested and has 5/5 Energy.\nLog in now to jump back into combat!"
			
			// Try to send email
			if userEmail != "" {
				err = w.sender.SendEmail(userEmail, subject, body)
				if err != nil {
					log.Printf("Failed to send energy email: %v", err)
				} else {
					log.Printf("Energy email sent successfully to %s", userEmail)
				}
			}

			// Add to notification history
			historyID := time.Now().Format("20060102150405") + "-" + userID
			_, _ = w.db.ExecContext(ctx, `
				INSERT INTO notification_history (id, user_id, type, subject, body)
				VALUES ($1, $2, $3, $4, $5)
			`, historyID, userID, "energy_refresh", subject, body)
		}
	}
}
