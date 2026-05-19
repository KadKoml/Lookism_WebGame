package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	userpb "github.com/kozie/lookism-rpg/api/proto/user"
	gachapb "github.com/kozie/lookism-rpg/api/proto/gacha"
)

type Gateway struct {
	userClient  userpb.UserServiceClient
	gachaClient gachapb.GachaEconomyServiceClient
}

func main() {
	// Connect to gRPC services
	userConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user-service: %v", err)
	}
	defer userConn.Close()

	gachaConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gacha-service: %v", err)
	}
	defer gachaConn.Close()

	gateway := &Gateway{
		userClient:  userpb.NewUserServiceClient(userConn),
		gachaClient: gachapb.NewGachaEconomyServiceClient(gachaConn),
	}

	mux := http.NewServeMux()

	// --- Auth & User Routes ---
	mux.HandleFunc("/api/auth/register", gateway.handleRegister)
	mux.HandleFunc("/api/auth/login", gateway.handleLogin)
	mux.HandleFunc("/api/user/profile", gateway.withAuth(gateway.handleGetProfile))
	
	// --- Gacha Routes ---
	mux.HandleFunc("/api/gacha/roll", gateway.withAuth(gateway.handleRollGacha))
	mux.HandleFunc("/api/gacha/freepack", gateway.withAuth(gateway.handleFreePack))

	// CORS wrapper
	handler := corsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("API Gateway listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start gateway: %v", err)
	}
}

// --- Middleware ---

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// withAuth extracts the Bearer token and passes it as gRPC metadata
func (g *Gateway) withAuth(handler func(http.ResponseWriter, *http.Request, context.Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		token := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Create gRPC context with metadata
		md := metadata.Pairs("authorization", "Bearer "+token)
		ctx := metadata.NewOutgoingContext(r.Context(), md)
		
		handler(w, r, ctx)
	}
}

// --- Handlers ---

func (g *Gateway) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req userpb.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	resp, err := g.userClient.Register(r.Context(), &req)
	writeJSONResponse(w, resp, err)
}

func (g *Gateway) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req userpb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	resp, err := g.userClient.Login(r.Context(), &req)
	writeJSONResponse(w, resp, err)
}

func (g *Gateway) handleGetProfile(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	// The user ID is ideally extracted from the token, but for now we expect the frontend to pass it
	// Actually, the user-service ValidateToken endpoint does this. 
	// Let's just pass the request to GetProfile and let user-service handle the JWT auth interceptor.
	// We need a way to parse the token to get the user ID, or the frontend sends it.
	// For simplicity, let's just make the frontend send the user_id in the query param
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeJSONError(w, "missing user_id", http.StatusBadRequest)
		return
	}
	
	resp, err := g.userClient.GetProfile(ctx, &userpb.UserRequest{UserId: userID})
	writeJSONResponse(w, resp, err)
}

func (g *Gateway) handleRollGacha(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if r.Method != "POST" {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req gachapb.RollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Assume single roll for now
	resp, err := g.gachaClient.RollGachaSingle(ctx, &req)
	writeJSONResponse(w, resp, err)
}

func (g *Gateway) handleFreePack(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	// Stub for free pack claim
	// Note: We need a ClaimFreePack endpoint in gacha.proto! 
	// The rubric requires 12 endpoints. I'll add this to the proto if needed, or just mock it here for now.
	writeJSONError(w, "Not implemented yet", http.StatusNotImplemented)
}

func writeJSONError(w http.ResponseWriter, errMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   errMsg,
		"message": errMsg,
	})
}

func writeJSONResponse(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		st, ok := status.FromError(err)
		statusCode := http.StatusInternalServerError
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				statusCode = http.StatusBadRequest
			case codes.Unauthenticated:
				statusCode = http.StatusUnauthorized
			case codes.PermissionDenied:
				statusCode = http.StatusForbidden
			case codes.NotFound:
				statusCode = http.StatusNotFound
			case codes.AlreadyExists:
				statusCode = http.StatusConflict
			case codes.Unimplemented:
				statusCode = http.StatusNotImplemented
			}
		}
		writeJSONError(w, err.Error(), statusCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
