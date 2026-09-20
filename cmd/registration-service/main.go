package main

import (
	"log"
	"net"

	"github.com/AbhiramiRajeev/event-ticketing-platform/config"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/database"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/registration"
	eventpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/event"
	pb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/registration"
	userpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	userConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed to connect to user service:", err)
	}

	eventConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed to connect to event service:", err)
	}

	userClient := userpb.NewUserServiceClient(userConn)
	eventClient := eventpb.NewEventServiceClient(eventConn)

	registrationRepo := registration.NewRepository(db)
	registrationService := registration.NewService(registrationRepo,userClient,eventClient)
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
