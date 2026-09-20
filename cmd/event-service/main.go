package main

import (
	"log"
	"net"

	"github.com/AbhiramiRajeev/event-ticketing-platform/config"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/database"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/event"
	pb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/event"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&event.Event{}); err != nil {
		log.Fatal("failed to migrate event table:", err)
	}

	eventRepo := event.NewRepository(db)
	eventService := event.NewService(eventRepo)
	grpcHandler := event.NewGRPCHandler(eventService)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterEventServiceServer(grpcServer, grpcHandler)

	log.Println("Event service running on :50052")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("failed to serve:", err)
	}
}