package registration

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/registration"
)

type GRPCHandler struct {
	pb.UnimplementedRegistrationServiceServer
	service *Service
}

func NewGRPCHandler(service *Service) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

func (h *GRPCHandler) CreateRegistration(
	ctx context.Context,
	req *pb.CreateRegistrationRequest,
) (*pb.Registration, error) {

	registration := &Registration{
		ID:     uuid.New().String(),
		UserID: req.UserId,
		EventID: req.EventId,
		Status: "registered",
	}

	err := h.service.CreateRegistration(registration,ctx)
	if err != nil {
		return nil, err
	}

	return &pb.Registration{
		Id:     registration.ID,
		UserId: registration.UserID,
		EventId: registration.EventID,
		Status: registration.Status,
	}, nil
}

func (h *GRPCHandler) GetRegistration(
	ctx context.Context,
	req *pb.GetRegistrationRequest,
) (*pb.Registration, error) {

	registration, err := h.service.GetRegistration(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Registration{
		Id:     registration.ID,
		UserId: registration.UserID,
		EventId: registration.EventID,
		Status: registration.Status,
	}, nil
}

func (h *GRPCHandler) GetRegistrations(
	ctx context.Context,
	req *pb.GetRegistrationsRequest,
) (*pb.GetRegistrationsResponse, error) {

	registrations, err := h.service.GetRegistrationsByUser(req.UserId)
	if err != nil {
		return nil, err
	}

	response := &pb.GetRegistrationsResponse{}

	for _, registration := range registrations {
		response.Registrations = append(
			response.Registrations,
			&pb.Registration{
				Id:     registration.ID,
				UserId: registration.UserID,
				EventId: registration.EventID,
				Status: registration.Status,
			},
		)
	}

	return response, nil
}

func (h *GRPCHandler) CancelRegistration(
	ctx context.Context,
	req *pb.CancelRegistrationRequest,
) (*pb.CancelRegistrationResponse, error) {

	err := h.service.CancelRegistration(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.CancelRegistrationResponse{
		Success: true,
	}, nil
}