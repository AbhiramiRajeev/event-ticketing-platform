package main

import (
	"log"
	"net"

	"github.com/AbhiramiRajeev/event-ticketing-platform/config"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/database"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/event"
	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/seats"
	pb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/event"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Migrate event and seat tables.
	if err := db.AutoMigrate(
		&event.Event{},
		&seats.Seat{},
	); err != nil {
		log.Fatal("failed to migrate event and seat tables:", err)
	}

	// Event dependencies.
	eventRepo := event.NewRepository(db)
	eventService := event.NewService(eventRepo)

	// Seat dependencies.
	seatRepo := seats.NewRepository(db)
	seatService := seats.NewService(seatRepo)

	// gRPC handler.
	grpcHandler := event.NewGRPCHandler(
		eventService,
		seatService,
	)

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