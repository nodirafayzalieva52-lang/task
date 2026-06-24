package service

import (
	"context"
	"errors"

	"project/internal/models"
	"project/internal/repository"
	"project/pkg/jwt"
	"project/pkg/password"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, request models.RegisterRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	exists, err := s.repo.ExistsByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("user with this email already exists")
	}

	passwordHash, err := password.Hash(request.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: passwordHash,
		Role:     models.UserRole,
	}

	err = s.repo.Add(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, request models.LoginRequest) (token string, err error) {
	err = request.Validate()
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByEmail(ctx, request.Email)
	if err != nil {
		return "", err
	}

	err = password.Compare(user.Password, request.Password)
	if err != nil {
		return "", err
	}

	token, err = jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil

}
