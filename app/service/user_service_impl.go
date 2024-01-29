package service

import (
	"context"

	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/entity"
	"github.com/Zainal21/my-ewallet/app/repositories"
)

type userServiceImpl struct {
	repo repositories.UserRepository
}

// CreateUser implements UserService.
func (u *userServiceImpl) CreateUser(ctx context.Context, payload dtos.UserRegistrationRequestDto) error {
	return u.repo.CreateUser(ctx, payload)
}

// GetUserByFieldName implements UserService.
func (u *userServiceImpl) GetUserByFieldName(ctx context.Context, fieldName string, value string) (*entity.User, error) {
	return u.repo.GetUserByFieldName(ctx, fieldName, value)
}

func NewUserServiceImpl(
	repo repositories.UserRepository,
) UserService {
	return &userServiceImpl{
		repo: repo,
	}
}
