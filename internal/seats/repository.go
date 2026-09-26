package seats

import (
	"errors"

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

func (r *Repository) Create(seat *Seat) error {
	return r.db.Create(seat).Error
}

func (r *Repository) CreateMany(seats []Seat) error {
	return r.db.Create(&seats).Error
}

func (r *Repository) GetByID(id string) (*Seat, error) {
	var seat Seat

	result := r.db.First(&seat, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &seat, nil
}

func (r *Repository) GetByEventID(eventID string) ([]Seat, error) {
	var seats []Seat

	result := r.db.Where("event_id = ?", eventID).Find(&seats)
	if result.Error != nil {
		return nil, result.Error
	}

	return seats, nil
}

func (r *Repository) GetAvailableSeats(eventID string) ([]Seat, error) {
	var seats []Seat

	result := r.db.
		Where("event_id = ? AND status = ?", eventID, "available").
		Find(&seats)

	if result.Error != nil {
		return nil, result.Error
	}

	return seats, nil
}

func (r *Repository) ReserveSeat(seatID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var seat Seat

		result := tx.First(&seat, "id = ?", seatID)

		if result.Error != nil {
			return result.Error
		}

		if seat.Status != "available" {
			return errors.New("seat is not available")
		}

		seat.Status = "reserved"

		return tx.Save(&seat).Error
	})
}

func (r *Repository) ReleaseSeat(seatID string) error {
	var seat Seat

	result := r.db.First(&seat, "id = ?", seatID)
	if result.Error != nil {
		return result.Error
	}

	if seat.Status != "reserved" {
		return errors.New("seat is not reserved")
	}

	seat.Status = "available"

	return r.db.Save(&seat).Error
}
