package seats

type Seat struct {
	ID         string `gorm:"primaryKey"`
	EventID    string `gorm:"uniqueIndex:idx_event_seat"`
	SeatNumber string `gorm:"uniqueIndex:idx_event_seat"`
	Status     string
}