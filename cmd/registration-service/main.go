package main

import (
	"log"
	"net"

	"github.com/AbhiramiRajeev/event-ticketing-platform/config"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/database"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/registration"
	pb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/registration"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&registration.Registration{}); err != nil {
		log.Fatal("failed to migrate registration table:", err)
	}

	registrationRepo := registration.NewRepository(db)
	registrationService := registration.NewService(registrationRepo)
	grpcHandler := registration.NewGRPCHandler(registrationService)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterRegistrationServiceServer(
		grpcServer,
		grpcHandler,
	)

	log.Println("Registration service running on :50053")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("failed to serve:", err)
	}
}