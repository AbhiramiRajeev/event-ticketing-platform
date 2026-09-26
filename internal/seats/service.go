package seats

import (
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetSeat(id string) (*Seat, error) {
	return s.repository.GetByID(id)
}

func (s *Service) GetSeatsByEvent(eventID string) ([]Seat, error) {
	return s.repository.GetByEventID(eventID)
}

func (s *Service) GetAvailableSeats(eventID string) ([]Seat, error) {
	return s.repository.GetAvailableSeats(eventID)
}

func (s *Service) ReserveSeat(seatID string) error {
	return s.repository.ReserveSeat(seatID)
}

func (s *Service) ReleaseSeat(seatID string) error {
	return s.repository.ReleaseSeat(seatID)
}

func (s *Service) CreateSeats(eventID string, capacity int) error {
	seats := make([]Seat, 0, capacity)

	for i := 1; i <= capacity; i++ {
		seats = append(seats, Seat{
			ID:         uuid.New().String(),
			EventID:    eventID,
			SeatNumber: fmt.Sprintf("A%d", i),
			Status:     "available",
		})
	}

	return s.repository.CreateMany(seats)
}
