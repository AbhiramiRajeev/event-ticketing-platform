package registration

type Registration struct {
	ID     string `gorm:"primaryKey"`
	UserID string
	EventID string
	Status string
}