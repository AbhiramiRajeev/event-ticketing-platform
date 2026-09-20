package user

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	pb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/user"
)

type GRPCHandler struct {
	pb.UnimplementedUserServiceServer
	service *Service
}

func NewGRPCHandler(service *Service) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

func (h *GRPCHandler) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.User, error) {

	user := &User{
		ID:    uuid.New().String(),
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	err := h.service.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (h *GRPCHandler) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.User, error) {

	user, err := h.service.GetUser(req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (h *GRPCHandler) GetUsers(
	ctx context.Context,
	req *pb.GetUsersRequest,
) (*pb.GetUsersResponse, error) {

	users, err := h.service.GetUsers()
	if err != nil {
		return nil, err
	}

	response := &pb.GetUsersResponse{}

	for _, user := range users {
		response.Users = append(response.Users, &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		})
	}

	return response, nil
}
