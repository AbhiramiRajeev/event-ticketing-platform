package main

import (
	"log"
	"net"

	"github.com/AbhiramiRajeev/event-ticketing-platform/config"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/database"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/user"
	pb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/user"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatal("failed to migrate user table:", err)
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	grpcHandler := user.NewGRPCHandler(userService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, grpcHandler)

	log.Println("User service running on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("failed to serve:", err)
	}
}