package user

type User struct {
	ID    string `gorm:"primaryKey"`
	Name  string
	Email string `gorm:"uniqueIndex"`
	Role  string
}