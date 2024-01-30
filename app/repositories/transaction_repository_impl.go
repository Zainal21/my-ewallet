package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/entity"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/utils/query"
	"github.com/Zainal21/my-ewallet/pkg/database/mysql"
	"github.com/google/uuid"
)

type transactionRepositoryImpl struct {
	db mysql.Adapter
}

// UpdateStatusTransactionLog implements TransactionRepository.
func (r *transactionRepositoryImpl) UpdateStatusTransactionLog(ctx context.Context, status, orderId string) error {
	_query := "UPDATE transactions SET status  = ? WHERE order_id = ?"
	_, err := r.db.Exec(ctx, _query, status, orderId)

	if err != nil {
		return err
	}

	return nil
}

const (
	pageSize   = 0
	dateFormat = "2006-01-02"
)

// BeginTx implements TransactionRepository.
func (t *transactionRepositoryImpl) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return t.db.BeginTx(ctx, opts)
}

// CreateDepositLog implements TransactionRepository.
func (t *transactionRepositoryImpl) CreateTransactionLog(ctx context.Context, payload dtos.TransactionDto) error {
	timeStr := helpers.GetTimeStrNow()
	Uuid := uuid.NewString()

	if _, err := t.db.Exec(ctx,
		`INSERT INTO transactions 
			(
				id, 
				order_id, 
				ref_id, 
				user_id,
				type, 
				gross_amount, 
				piece, 
				amount, 
				note, 
				status, 
				created_at, 
				updated_at
			) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`,
		Uuid,
		payload.OrderID,
		payload.RefID,
		payload.UserID,
		payload.Type,
		payload.GrossAmount,
		payload.Piece,
		payload.Amount,
		payload.Note,
		payload.Status,
		&timeStr,
		&timeStr,
	); err != nil {
		return err
	}

	return nil
}

// CreateDepositLog implements TransactionRepository.
func (t *transactionRepositoryImpl) CreateDepositLog(ctx context.Context, payload dtos.LedgerDto) error {
	timeStr := helpers.GetTimeStrNow()
	Uuid := uuid.NewString()

	if _, err := t.db.Exec(ctx,
		`INSERT INTO ledgers 
			(
				id, 
				user_id, 
				ref_id,
				type, 
				current_deposit, 
				change_deposit, 
				final_deposit, 
				note, 
				created_at, 
				updated_at
			) 
		VALUES (?,?,?,?,?,?,?,?,?,?)`,
		Uuid,
		payload.UserID,
		payload.RefID,
		payload.Type,
		payload.CurrentDeposit,
		payload.ChangeDeposit,
		payload.FinalDeposit,
		payload.Note,
		&timeStr,
		&timeStr,
	); err != nil {
		return err
	}

	return nil
}

// GetBalance implements TransactionRepository.
func (t *transactionRepositoryImpl) GetTransactionByFieldName(ctx context.Context, fieldName string, value string) (*entity.Transaction, error) {
	transQuery := query.SelectQuery(
		"transactions",
		[]string{
			"id",
			"user_id",
			"ref_id",
			"type",
			"gross_amount",
			"amount",
			"note",
			"created_at",
			"updated_at",
		},
		fieldName+" = ? ORDER BY created_at DESC",
		1,
		0,
	)

	var result entity.Transaction

	row := t.db.QueryRowX(ctx, transQuery, value)

	if err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.RefID,
		&result.Type,
		&result.GrossAmount,
		&result.Amount,
		&result.Note,
		&result.CreatedAt,
		&result.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBalance implements TransactionRepository.
func (t *transactionRepositoryImpl) GetBalance(ctx context.Context, fieldName string, value string) (*entity.Ledger, error) {
	balanceQuery := query.SelectQuery(
		"ledgers",
		[]string{
			"id",
			"user_id",
			"ref_id",
			"type",
			"current_deposit",
			"change_deposit",
			"final_deposit",
			"note",
			"created_at",
			"updated_at",
		},
		fieldName+" = ? ORDER BY created_at DESC",
		1,
		0,
	)

	var result entity.Ledger

	row := t.db.QueryRowX(ctx, balanceQuery, value)

	err := row.Scan(
		&result.Id,
		&result.UserID,
		&result.RefID,
		&result.Type,
		&result.CurrentDeposit,
		&result.ChangeDeposit,
		&result.FinalDeposit,
		&result.Note,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDepositHistory implements TransactionRepository.
func (t *transactionRepositoryImpl) GetDepositHistory(ctx context.Context, payload dtos.TransactionRequestDto) (*[]entity.Ledger, int, error) {
	offset := (payload.Page - 1) * pageSize

	countQuery := "SELECT 1 as record FROM ledgers"
	var totalCount int = 0
	countQueryFinal := `SELECT COALESCE(SUM(record), 0) AS total FROM (` + countQuery + `) AS subquery`
	err := t.db.QueryRowX(ctx, countQueryFinal).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	transQuery := `SELECT 
			L.id, 
			L.user_id, 
			L.ref_id, 
			L.type, 
			L.current_deposit, 
			L.change_deposit, 
			L.final_deposit, 
			L.note, 
			L.created_at, 
			L.updated_at
		FROM ledgers as L
		JOIN transactions AS T
			ON L.id = T.transaction_id
		WHERE user_id = ?`

	transQuery = t.searchLogReportFilter(transQuery, payload.Search)
	transQuery = t.dateFilterLogReport(transQuery, payload.DateFrom, payload.DateTo)

	transQuery += " ORDER BY created_at DESC LIMIT 10 OFFSET ?"

	var result []entity.Ledger
	err = t.db.Query(ctx, &result, transQuery, payload.UserId, offset)

	if err != nil {
		return nil, 0, err
	}

	return &result, totalCount, err
}

func (l *transactionRepositoryImpl) searchLogReportFilter(_query, search string) string {
	if search != "" {
		searchQuery := fmt.Sprintf(`
				AND (
					LOWER(type) LIKE '%%%s%%' OR
					LOWER(note) LIKE '%%%s%%'
				)`, search, search)
		_query += searchQuery
	}
	return _query
}

func (l *transactionRepositoryImpl) dateFilterLogReport(_query, startDate, endDate string) string {
	if startDate != "" && endDate != "" {
		startDateTime, _ := time.Parse("02-01-2006", startDate)
		endDateTime, _ := time.Parse("02-01-2006", endDate)
		filterDateQuery := fmt.Sprintf(`
        AND DATE(ledgers.created_at) BETWEEN '%s' AND '%s'`, startDateTime.Format(dateFormat), endDateTime.Format(dateFormat))
		_query += filterDateQuery
	}
	return _query
}

func NewTransactionRepositoryImpl(db mysql.Adapter) TransactionRepository {
	return &transactionRepositoryImpl{
		db: db,
	}
}
