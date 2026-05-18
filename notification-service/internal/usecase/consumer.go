package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	"github.com/nats-io/nats.go"
)

type FreePackEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Timestamp string `json:"timestamp"`
}

type NotificationUsecase struct {
	nc       *nats.Conn
	smtpHost string
	smtpPort string
}

func NewNotificationUsecase(nc *nats.Conn, host, port string) *NotificationUsecase {
	return &NotificationUsecase{
		nc:       nc,
		smtpHost: host, // Default: localhost or mailpit
		smtpPort: port, // Default: 1025
	}
}

func (u *NotificationUsecase) StartConsumer() {
	_, err := u.nc.Subscribe("notifications.free_pack", func(m *nats.Msg) {
		var event FreePackEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Printf("Failed to unmarshal FreePackEvent: %v", err)
			return
		}

		log.Printf("Received FreePackReady event for user: %s (Email: %s)", event.UserID, event.Email)
		
		err = u.SendEmail(event.Email, "Your Free Character Pack is Ready!", "Log in now to claim your free Lookism character pack!")
		if err != nil {
			log.Printf("Failed to send email to %s: %v", event.Email, err)
		} else {
			log.Printf("Successfully sent free pack email to %s", event.Email)
		}
	})

	if err != nil {
		log.Fatalf("Failed to subscribe to notifications.free_pack: %v", err)
	}
}

func (u *NotificationUsecase) SendEmail(to, subject, body string) error {
	from := "noreply@lookism-rpg.local"
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	// Assuming unauthenticated local SMTP (Mailpit)
	addr := fmt.Sprintf("%s:%s", u.smtpHost, u.smtpPort)
	return smtp.SendMail(addr, nil, from, []string{to}, msg)
}
