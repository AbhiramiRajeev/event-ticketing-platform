package registration

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateRegistration(registration *Registration) error {
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