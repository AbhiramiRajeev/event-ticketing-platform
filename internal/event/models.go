package event


type Event struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	Venue       string
	Date        string
}
