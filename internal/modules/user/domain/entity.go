package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID       int64  `gorm:"primaryKey"`
	Name     string `gorm:"not null;type:varchar(255)"`
	Email    string `gorm:"not null;unique;type:varchar(255)"`
	Password string `json:"-" gorm:"not null"`
}

type UserRepository interface {
	GetByID(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(user *User) error
}
