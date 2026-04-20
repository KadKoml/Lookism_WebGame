package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	pb "github.com/kozie/lookism-rpg/api/proto/user"
	delivery "github.com/kozie/lookism-rpg/user-service/internal/delivery/grpc"
	"github.com/kozie/lookism-rpg/user-service/internal/repository/postgres"
	"github.com/kozie/lookism-rpg/user-service/internal/usecase"
)

func main() {
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://lookism_user:lookism_password@localhost:5432/lookism_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Run migrations (simple approach — execute the SQL file)
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

	// Initialize repositories (Clean Architecture: depend on interfaces)
	userRepo := postgres.NewUserRepository(db)
	cardRepo := postgres.NewCardRepository(db)
	squadRepo := postgres.NewSquadRepository(db)

	// Initialize usecase
	uc := usecase.NewUserUsecase(userRepo, cardRepo, squadRepo)

	// Initialize gRPC handler
	handler := delivery.NewUserHandler(uc)

	// Start gRPC server
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, handler)

	go func() {
		log.Printf("User Service gRPC server listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down User Service...")
	grpcServer.GracefulStop()
}
