package event

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateEvent(event *Event) error {
	return s.repository.Create(event)
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