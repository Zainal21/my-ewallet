package repositories

import (
	"context"
	"database/sql"

	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/entity"
)

type UserRepository interface {
	GetUserByFieldName(ctx context.Context, fieldName string, value string) (*entity.User, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type PersonalTokenRepository interface {
	// create token
	Create(ctx context.Context, personalTokenDto *dtos.PersonalAccessTokenDto) (string, error)
	// verify token
	Verify(ctx context.Context, token string) (*entity.User, error)
	//delete token by token
	Delete(ctx context.Context, token string) error
	// delete token by user id
	DeleteByUserId(ctx context.Context, user_id string) error
}
