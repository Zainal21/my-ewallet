package service

import (
	"context"

	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/entity"
)

type UserService interface {
	// get spesific user by field name
	GetUserByFieldName(ctx context.Context, fieldName, value string) (*entity.User, error)
	// create user data
	CreateUser(ctx context.Context, payload dtos.UserRegistrationRequestDto) error
}

type TransactionService interface {
	// get spesific user by field name
	GetBalance(ctx context.Context, fieldName string, value string) (*entity.Ledger, error)
	// get transaction history
	GetDepositHistory(ctx context.Context, payload dtos.TransactionRequestDto) (*[]entity.Ledger, int, error)
	// create transaction (topup, transfer/payment)
	CreateDepositLog(ctx context.Context, payload dtos.LedgerDto) error
	// craete transaction log
	CreateTransactionLog(ctx context.Context, payload dtos.TransactionDto) error
	// update status transaction log
	UpdateStatusTransactionLog(ctx context.Context, status, orderId string) error
	// get transcation by fieldname
	GetTransactionByFieldName(ctx context.Context, fieldName string, value string) (*entity.Transaction, error)
}
