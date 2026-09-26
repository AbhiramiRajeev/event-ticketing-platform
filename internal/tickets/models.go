package tickets

type Ticket struct {
	ID             string `gorm:"primaryKey"`
	RegistrationID string `gorm:"not null"`
	SeatID         string `gorm:"not null"`
	Status         string
}