package consumer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/kozie/lookism-rpg/notification-service/internal/service"
)

type FreePackEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

type NotificationConsumer struct {
	nc     *nats.Conn
	sender service.SMTPSender
}

func NewNotificationConsumer(nc *nats.Conn, sender service.SMTPSender) *NotificationConsumer {
	return &NotificationConsumer{
		nc:     nc,
		sender: sender,
	}
}

func (c *NotificationConsumer) Start() error {
	_, err := c.nc.Subscribe("notifications.free_pack", func(msg *nats.Msg) {
		var event FreePackEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			return
		}

		log.Printf("Received free pack event for user %s (email: %s)", event.UserID, event.Email)

		// Hardcoded fallback if email is empty for testing
		targetEmail := event.Email
		if targetEmail == "" {
			targetEmail = "player@lookism.local" 
		}

		subject := "Your Free Pack is Ready! 🎉"
		body := "Hello!\n\nYour free gacha pack in Lookism RPG is now ready to claim.\nLog in now to see who you get!\n\n- The Lookism Team"

		err := c.sender.SendEmail(targetEmail, subject, body)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", targetEmail, err)
		} else {
			log.Printf("Successfully sent email to %s", targetEmail)
		}
	})

	if err != nil {
		return err
	}

	log.Println("Notification consumer started, listening on 'notifications.free_pack'")
	return nil
}
