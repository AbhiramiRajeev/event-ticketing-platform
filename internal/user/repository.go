package user

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetByID(id string) (*User, error) {
	var user User

	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *Repository) GetAll() ([]User, error) {
	var users []User

	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}	

	return users, nil
}