package repository

import (
	"errors"
	"fmt"

	"github.com/geraldiaditya/ratix-backend/internal/modules/user/domain"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	DB *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) GetByID(id int64) (*domain.User, error) {
	var user domain.User
	if err := r.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}
