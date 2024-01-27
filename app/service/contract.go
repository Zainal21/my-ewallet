package service

import (
	"context"

	"github.com/Zainal21/my-ewallet/app/entity"
)

type UserService interface {
	GetUserByFieldName(ctx context.Context, fieldName, value string) (*entity.User, error)
}
