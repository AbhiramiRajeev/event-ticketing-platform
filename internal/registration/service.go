package registration

import (
	"context"

	eventpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/event"
	userpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/user"
)

type Service struct {
	repository  *Repository
	userClient  userpb.UserServiceClient
	eventClient eventpb.EventServiceClient
}

func NewService(
	repository *Repository,
	userClient userpb.UserServiceClient,
	eventClient eventpb.EventServiceClient,
) *Service {
	return &Service{
		repository:  repository,
		userClient:  userClient,
		eventClient: eventClient,
	}
}

func (s *Service) CreateRegistration(registration *Registration,ctx context.Context) error {

	_, err := s.userClient.GetUser(
		ctx,
		&userpb.GetUserRequest{
			Id: registration.UserID,
		},
	)
	if err != nil {
		return err
	}

	_, err = s.eventClient.GetEvent(
		ctx,
		&eventpb.GetEventRequest{
			Id: registration.EventID,
		},
	)
	if err != nil {
		return err
	}

	return s.repository.Create(registration)
}

func (s *Service) GetRegistration(id string) (*Registration, error) {
	return s.repository.Get(id)
}

func (s *Service) GetRegistrationsByUser(userID string) ([]Registration, error) {
	return s.repository.GetByUserID(userID)
}

func (s *Service) CancelRegistration(id string) error {
	registration, err := s.repository.Get(id)
	if err != nil {
		return err
	}

	registration.Status = "Cancelled"

	return s.repository.Update(registration)
}
