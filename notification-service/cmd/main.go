package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	pb "github.com/kozie/lookism-rpg/api/proto/notification"
	"github.com/kozie/lookism-rpg/notification-service/internal/consumer"
	delivery "github.com/kozie/lookism-rpg/notification-service/internal/delivery/grpc"
	"github.com/kozie/lookism-rpg/notification-service/internal/service"
)

func main() {
	// Connect to PostgreSQL (needed for EnergyWorker and Notification History)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://lookism_user:lookism_password@localhost:5432/lookism_db?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	// Run migrations
	migrationSQL, err := os.ReadFile("migrations/001_init.up.sql")
	if err != nil {
		log.Printf("Warning: Could not read migration file: %v", err)
	} else {
		if _, err := db.Exec(string(migrationSQL)); err != nil {
			log.Printf("Warning: Migration may have already been applied: %v", err)
		} else {
			log.Println("Migrations applied successfully")
		}
	}

	// Connect to NATS
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()
	log.Println("Connected to NATS")

	// Setup SMTP Sender
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com" // Default to Gmail for final project
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587" // Default to Gmail TLS port
	}
	
	sender := service.NewSMTPSender(smtpHost, smtpPort)

	// Start Background Consumers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notifConsumer := consumer.NewNotificationConsumer(nc, sender)
	if err := notifConsumer.Start(); err != nil {
		log.Fatalf("Failed to start NATS consumer: %v", err)
	}

	energyWorker := consumer.NewEnergyWorker(db, sender)
	go energyWorker.Start(ctx)

	// Setup gRPC Server
	handler := delivery.NewNotificationHandler()
	port := os.Getenv("PORT")
	if port == "" {
		port = "50054"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, handler)

	go func() {
		log.Printf("Notification Service gRPC server listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Keep the service running
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Notification Service...")
	cancel() // Stop energy worker
	grpcServer.GracefulStop()
}
