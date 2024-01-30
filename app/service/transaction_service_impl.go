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

// GetTransactionByFieldName implements TransactionService.
func (t *transactionServiceImpl) GetTransactionByFieldName(ctx context.Context, fieldName string, value string) (*entity.Transaction, error) {
	return t.transRepo.GetTransactionByFieldName(ctx, fieldName, value)
}

// UpdateStatusTransactionLog implements TransactionService.
func (t *transactionServiceImpl) UpdateStatusTransactionLog(ctx context.Context, status string, refId string) error {
	return t.transRepo.UpdateStatusTransactionLog(ctx, status, refId)
}

// CreateTransactionLog implements TransactionService.
func (t *transactionServiceImpl) CreateTransactionLog(ctx context.Context, payload dtos.TransactionDto) error {
	return t.transRepo.CreateTransactionLog(ctx, payload)
}

// CreateDepositLog implements TransactionService.
func (t *transactionServiceImpl) CreateDepositLog(ctx context.Context, payload dtos.LedgerDto) error {
	return t.transRepo.CreateDepositLog(ctx, payload)
}

// GetBalance implements TransactionService.
func (t *transactionServiceImpl) GetBalance(ctx context.Context, fieldName string, value string) (*entity.Ledger, error) {
	return t.transRepo.GetBalance(ctx, fieldName, value)
}

// GetDepositHistory implements TransactionService.
func (t *transactionServiceImpl) GetDepositHistory(ctx context.Context, payload dtos.TransactionRequestDto) (*[]entity.Ledger, int, error) {
	return t.transRepo.GetDepositHistory(ctx, payload)
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
