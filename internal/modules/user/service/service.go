package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/geraldiaditya/ratix-backend/internal/modules/user/domain"
	"github.com/geraldiaditya/ratix-backend/internal/modules/user/dto"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo      domain.UserRepository
	JWTSecret string
}

func NewUserService(repo domain.UserRepository, jwtSecret string) *UserService {
	return &UserService{Repo: repo, JWTSecret: jwtSecret}
}

func (s *UserService) GetUser(id int64) (*domain.User, error) {
	return s.Repo.GetByID(id)
}

func (s *UserService) RegisterUser(name, email, password, confirmPassword string) (*domain.User, error) {
	if password != confirmPassword {
		return nil, errors.New("passwords do not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	err = s.Repo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(email, password string) (*dto.LoginResponse, error) {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	t, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: t,
		User:  dto.ToUserResponse(user),
	}, nil
}
