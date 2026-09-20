package registration

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(reg *Registration) error {
	return r.db.Create(reg).Error
}

func (r *Repository) Get(id string) (*Registration, error) {

	var reg Registration
	result := r.db.First(&reg, "id=?", id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &reg, nil
}

func (r *Repository) GetByUserID(userID string) ([]Registration, error) {
	var registrations []Registration

	result := r.db.Where("user_id = ?", userID).Find(&registrations)
	if result.Error != nil {
		return nil, result.Error
	}

	return registrations, nil
}

func (r *Repository) Update(registration *Registration) error {
	return r.db.Save(registration).Error
}