package event

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/AbhiramiRajeev/event-ticketing-platform/internal/seats"
)

type Service struct {
	repository     *Repository
	seatRepository *seats.Repository
	db             *gorm.DB
}

func NewService(
	db *gorm.DB,
	repository *Repository,
	seatRepository *seats.Repository,
) *Service {
	return &Service{
		db:             db,
		repository:     repository,
		seatRepository: seatRepository,
	}
}

func (s *Service) CreateEvent(event *Event) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		if err := s.repository.CreateWithDB(tx, event); err != nil {
			return err
		}

		seats := make([]seats.Seat, 0, event.Capacity)

		for i := 1; i <= event.Capacity; i++ {
			seats = append(seats, seats.Seat{
				ID:         uuid.New().String(),
				EventID:    event.ID,
				SeatNumber: fmt.Sprintf("A%d", i),
				Status:     "available",
			})
		}

		if err := s.seatRepository.CreateManyWithDB(tx, seats); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) GetEvent(id string) (*Event, error) {
	return s.repository.GetEvent(id)
}

func (s *Service) GetEvents() ([]Event, error) {
	return s.repository.GetAll()
}

func (s *Service) UpdateEvent(event *Event) error {
	return s.repository.Update(event)
}

func (s *Service) DeleteEvent(id string) error {
	return s.repository.Delete(id)
}