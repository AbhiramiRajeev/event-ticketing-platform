package event

import "time"

type Event struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	Venue       string
	Date        time.Time
	Capacity    int
}
