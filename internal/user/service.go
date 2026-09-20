package user

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateUser(user *User) error {
	return s.repository.Create(user)
}

func (s *Service) GetUser(id string) (*User, error) {
	return s.repository.GetByID(id)
}

func (s *Service) GetUsers() ([]User, error) {
	return s.repository.GetAll()
}