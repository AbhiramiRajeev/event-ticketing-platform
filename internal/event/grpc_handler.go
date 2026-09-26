package event

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	pb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/event"
)

type GRPCHandler struct {
    pb.UnimplementedEventServiceServer
    service      *Service
    seatService  *seats.Service
}

func NewGRPCHandler(
    service *Service,
    seatService *seats.Service,
) *GRPCHandler {
    return &GRPCHandler{
        service:     service,
        seatService: seatService,
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
		Date:        timestamppb.New(event.Date),
		Capacity:    int32(event.Capacity),
	}, nil
}

func (h *GRPCHandler) GetEvents(
	ctx context.Context,
	req *pb.GetEventsRequest,
) (*pb.GetEventsResponse, error) {

	events, err := h.service.GetEvents()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get events")
	}

	response := &pb.GetEventsResponse{
		Events: make([]*pb.Event, 0, len(events)),
	}

	for _, event := range events {
		response.Events = append(response.Events, &pb.Event{
			Id:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			Venue:       event.Venue,
			Date:        timestamppb.New(event.Date),
			Capacity:    int32(event.Capacity),
		})
	}

	return response, nil
}
func (s *Service) CreateEvent(event *Event) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		// Create the event.
		if err := s.repository.CreateWithDB(tx, event); err != nil {
			return err
		}

		// Generate seats for the event.
		seats := make([]seats.Seat, 0, event.Capacity)

		for i := 1; i <= event.Capacity; i++ {
			seats = append(seats, seats.Seat{
				ID:         uuid.New().String(),
				EventID:    event.ID,
				SeatNumber: fmt.Sprintf("A%d", i),
				Status:     "available",
			})
		}

		// Create all seats in the same transaction.
		if err := s.seatRepository.CreateManyWithDB(tx, seats); err != nil {
			return err
		}

		return nil
	})
}

func (h *GRPCHandler) UpdateEvent(
	ctx context.Context,
	req *pb.UpdateEventRequest,
) (*pb.Event, error) {

	if req.Date == nil {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}

	if req.Capacity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "capacity must be greater than 0")
	}

	event := &Event{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Venue:       req.Venue,
		Date:        req.Date.AsTime(),
		Capacity:    int(req.Capacity),
	}

	err := h.service.UpdateEvent(event)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}

		return nil, status.Error(codes.Internal, "failed to update event")
	}

	return &pb.Event{
		Id:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Venue:       event.Venue,
		Date:        timestamppb.New(event.Date),
		Capacity:    int32(event.Capacity),
	}, nil
}

func (h *GRPCHandler) DeleteEvent(
	ctx context.Context,
	req *pb.DeleteEventRequest,
) (*pb.DeleteEventResponse, error) {

	err := h.service.DeleteEvent(req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}

		return nil, status.Error(codes.Internal, "failed to delete event")
	}

	return &pb.DeleteEventResponse{
		Success: true,
	}, nil
}


func (h *GRPCHandler) GetSeats(
    ctx context.Context,
    req *pb.GetSeatsRequest,
) (*pb.GetSeatsResponse, error) {

    seats, err := h.seatService.GetSeatsByEvent(req.EventId)
    if err != nil {
        return nil, status.Error(
            codes.Internal,
            "failed to get seats",
        )
    }

    response := &pb.GetSeatsResponse{
        Seats: make([]*pb.Seat, 0, len(seats)),
    }

    for _, seat := range seats {
        response.Seats = append(response.Seats, &pb.Seat{
            Id:         seat.ID,
            EventId:    seat.EventID,
            SeatNumber: seat.SeatNumber,
            Status:     seat.Status,
        })
    }

    return response, nil
}

func (h *GRPCHandler) GetAvailableSeats(
	ctx context.Context,
	req *pb.GetAvailableSeatsRequest,
) (*pb.GetAvailableSeatsResponse, error) {

	seats, err := h.seatService.GetAvailableSeats(req.EventId)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to get available seats",
		)
	}

	response := &pb.GetAvailableSeatsResponse{
		Seats: make([]*pb.Seat, 0, len(seats)),
	}

	for _, seat := range seats {
		response.Seats = append(response.Seats, &pb.Seat{
			Id:         seat.ID,
			EventId:    seat.EventID,
			SeatNumber: seat.SeatNumber,
			Status:     seat.Status,
		})
	}

	return response, nil
}

func (h *GRPCHandler) ReserveSeat(
	ctx context.Context,
	req *pb.ReserveSeatRequest,
) (*pb.Seat, error) {

	err := h.seatService.ReserveSeat(req.SeatId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"seat not found",
			)
		}

		return nil, status.Error(
			codes.FailedPrecondition,
			err.Error(),
		)
	}

	seat, err := h.seatService.GetSeat(req.SeatId)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to get reserved seat",
		)
	}

	return &pb.Seat{
		Id:         seat.ID,
		EventId:    seat.EventID,
		SeatNumber: seat.SeatNumber,
		Status:     seat.Status,
	}, nil
}


func (h *GRPCHandler) ReleaseSeat(
	ctx context.Context,
	req *pb.ReleaseSeatRequest,
) (*pb.Seat, error) {

	err := h.seatService.ReleaseSeat(req.SeatId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"seat not found",
			)
		}

		return nil, status.Error(
			codes.FailedPrecondition,
			err.Error(),
		)
	}

	seat, err := h.seatService.GetSeat(req.SeatId)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to get released seat",
		)
	}

	return &pb.Seat{
		Id:         seat.ID,
		EventId:    seat.EventID,
		SeatNumber: seat.SeatNumber,
		Status:     seat.Status,
	}, nil
}