package event

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(event *Event) error {
	return r.db.Create(event).Error
}

func (r *Repository) GetEvent(id string) (*Event, error) {
	var event Event

	result := r.db.First(&event, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &event, nil
}

func (r *Repository) GetAll() ([]Event, error) {
	var events []Event

	result := r.db.Find(&events)

	if result.Error != nil {
		return nil, result.Error
	}

	return events, nil
}

func (r *Repository) Update(event *Event) error {
	return r.db.Save(event).Error
}

func (r *Repository) Delete(id string) error {
	result := r.db.Delete(&Event{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}