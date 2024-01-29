package service

import (
	"context"

	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/entity"
	"github.com/Zainal21/my-ewallet/app/repositories"
)

type transactionServiceImpl struct {
	repo      repositories.UserRepository
	transRepo repositories.TransactionRepository
}

// CreateTransaction implements TransactionService.
func (t *transactionServiceImpl) CreateTransaction(ctx context.Context, payload dtos.LedgerDto) error {
	return t.transRepo.CreateTransaction(ctx, payload)
}

// GetBalance implements TransactionService.
func (t *transactionServiceImpl) GetBalance(ctx context.Context, fieldName string, value string) (*entity.Ledger, error) {
	return t.transRepo.GetBalance(ctx, fieldName, value)
}

// GetTransactionHistory implements TransactionService.
func (t *transactionServiceImpl) GetTransactionHistory(ctx context.Context, payload dtos.TransactionRequestDto) (*[]entity.Ledger, int, error) {
	return t.transRepo.GetTransactionHistory(ctx, payload)
}

func NewTransactionServiceImpl(
	repo repositories.UserRepository,
	transRepo repositories.TransactionRepository,
) TransactionService {
	return &transactionServiceImpl{
		repo:      repo,
		transRepo: transRepo,
	}
}
