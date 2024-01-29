package repositories

import (
	"context"
	"database/sql"

	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/entity"
)

type UserRepository interface {
	// get spesific user by field name
	GetUserByFieldName(ctx context.Context, fieldName string, value string) (*entity.User, error)
	// begin database transaction
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	// create users
	CreateUser(ctx context.Context, payload dtos.UserRegistrationRequestDto) error
}

type TransactionRepository interface {
	// get spesific user by field name
	GetBalance(ctx context.Context, fieldName string, value string) (*entity.Ledger, error)
	// get transaction history
	GetTransactionHistory(ctx context.Context, payload dtos.TransactionRequestDto) (*[]entity.Ledger, int, error)
	// begin database transaction
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	// create transaction (topup, transfer/payment)
	CreateTransaction(ctx context.Context, payload dtos.LedgerDto) error
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
