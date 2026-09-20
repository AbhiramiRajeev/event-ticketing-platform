package main

import (
	"log"

	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/gateway"
	eventpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/event"

	// registrationpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/registration"
	// userpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/user"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {

	// User Service
	userConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed to connect to user service:", err)
	}
	defer userConn.Close()

	// Event Service
	eventConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed to connect to event service:", err)
	}
	defer eventConn.Close()

	// Registration Service
	registrationConn, err := grpc.NewClient(
		"localhost:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed to connect to registration service:", err)
	}
	defer registrationConn.Close()

	// Create gRPC clients
	// userClient := userpb.NewUserServiceClient(userConn)
	eventClient := eventpb.NewEventServiceClient(eventConn)
	// registrationClient := registrationpb.NewRegistrationServiceClient(registrationConn)

	// Create HTTP handlers
	eventHandler := gateway.NewEventHandler(eventClient)

	// Gin router
	r := gin.Default()

	// Event routes
	r.GET("/events", eventHandler.GetEvents)
	r.GET("/events/:id", eventHandler.GetEvent)
	r.POST("/events", eventHandler.CreateEvent)
	r.PUT("/events/:id", eventHandler.UpdateEvent)
	r.DELETE("/events/:id", eventHandler.DeleteEvent)

	log.Println("API Gateway running on :8080")

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
