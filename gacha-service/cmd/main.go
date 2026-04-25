package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/kozie/lookism-rpg/api/proto/gacha"
	delivery "github.com/kozie/lookism-rpg/gacha-service/internal/delivery/grpc"
	"github.com/kozie/lookism-rpg/gacha-service/internal/repository/postgres"
	redisrepo "github.com/kozie/lookism-rpg/gacha-service/internal/repository/redis"
	"github.com/kozie/lookism-rpg/gacha-service/internal/usecase"
)

func main() {
	// Connect to PostgreSQL
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

	// Connect to Redis
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer rdb.Close()

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

	// Connect to User Service
	userSvcAddr := os.Getenv("USER_SERVICE_ADDR")
	if userSvcAddr == "" {
		userSvcAddr = "localhost:50051"
	}
	userConn, err := grpc.NewClient(userSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user-service: %v", err)
	}
	defer userConn.Close()

	// Repositories
	economyRepo := postgres.NewEconomyRepository(db)
	questRepo := postgres.NewQuestRepository(db)
	cache := redisrepo.NewGachaCache(rdb)

	// Usecase
	uc := usecase.NewGachaUsecase(economyRepo, questRepo, cache, nc, userConn)

	// gRPC Handler
	handler := delivery.NewGachaHandler(uc)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "50052"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGachaEconomyServiceServer(grpcServer, handler)

	go func() {
		log.Printf("Gacha Service gRPC server listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Gacha Service...")
	grpcServer.GracefulStop()
}
