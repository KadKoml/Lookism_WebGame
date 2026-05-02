package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/kozie/lookism-rpg/api/proto/combat"
	delivery "github.com/kozie/lookism-rpg/combat-service/internal/delivery/grpc"
	redisrepo "github.com/kozie/lookism-rpg/combat-service/internal/repository/redis"
	"github.com/kozie/lookism-rpg/combat-service/internal/usecase"
)

func main() {
	// Connect to Redis
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer rdb.Close()
	
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

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
	combatRepo := redisrepo.NewCombatRepository(rdb)

	// Usecase
	uc := usecase.NewCombatUsecase(combatRepo, nc, userConn)

	// gRPC Handler
	handler := delivery.NewCombatHandler(uc)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "50053"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCombatEngineServiceServer(grpcServer, handler)

	go func() {
		log.Printf("Combat Service gRPC server listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Combat Service...")
	grpcServer.GracefulStop()
}
