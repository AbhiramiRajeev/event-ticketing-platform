package event

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	pb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/event"
)

type GRPCHandler struct {
	pb.UnimplementedEventServiceServer
	service *Service
}

func NewGRPCHandler(service *Service) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

func (h *GRPCHandler) GetEvent(
	ctx context.Context,
	req *pb.GetEventRequest,
) (*pb.Event, error) {

	event, err := h.service.GetEvent(req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}

		return nil, status.Error(codes.Internal, "failed to get event")
	}

	return &pb.Event{
		Id:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Venue:       event.Venue,
		Date:        event.Date,
	}, nil
}

func (h *GRPCHandler) GetEvents(
	ctx context.Context,
	req *pb.GetEventsRequest,
) (*pb.GetEventsResponse, error) {

	events, err := h.service.GetEvents()
	if err != nil {
		return nil, err
	}

	response := &pb.GetEventsResponse{}

	for _, event := range events {
		response.Events = append(response.Events, &pb.Event{
			Id:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			Venue:       event.Venue,
			Date:        event.Date,
		})
	}

	return response, nil
}

func (h *GRPCHandler) CreateEvent(
	ctx context.Context,
	req *pb.CreateEventRequest,
) (*pb.Event, error) {

	event := &Event{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Venue:       req.Venue,
		Date:        req.Date,
	}

	err := h.service.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	return &pb.Event{
		Id:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Venue:       event.Venue,
		Date:        event.Date,
	}, nil
}
func (h *GRPCHandler) UpdateEvent(
	ctx context.Context,
	req *pb.UpdateEventRequest,
) (*pb.Event, error) {

	event := &Event{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Venue:       req.Venue,
		Date:        req.Date,
	}

	err := h.service.UpdateEvent(event)
	if err != nil {
		return nil, err
	}

	return &pb.Event{
		Id:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Venue:       event.Venue,
		Date:        event.Date,
	}, nil
}

func (h *GRPCHandler) DeleteEvent(
	ctx context.Context,
	req *pb.DeleteEventRequest,
) (*pb.DeleteEventResponse, error) {

	err := h.service.DeleteEvent(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteEventResponse{
		Success: true,
	}, nil
}
