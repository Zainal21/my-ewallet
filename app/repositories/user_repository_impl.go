package repositories

import (
	"context"
	"database/sql"

	"github.com/Zainal21/my-ewallet/app/dtos"
	"github.com/Zainal21/my-ewallet/app/entity"
	"github.com/Zainal21/my-ewallet/app/helpers"
	"github.com/Zainal21/my-ewallet/app/utils/query"
	"github.com/Zainal21/my-ewallet/pkg/database/mysql"
	"github.com/google/uuid"
)

type userRepositoryImpl struct {
	db mysql.Adapter
}

// CreateUser implements UserRepository.
func (r *userRepositoryImpl) CreateUser(ctx context.Context, payload dtos.UserRegistrationRequestDto) error {
	timeStr := helpers.GetTimeStrNow()
	Uuid := uuid.NewString()
	passwordHash, _ := helpers.HashPassword(payload.Password)

	if _, err := r.db.Exec(ctx,
		`INSERT INTO users 
			(
				id, 
				name, 
				email,
				phone_number, 
				password, 
				created_at, 
				updated_at
			) 
		VALUES (?, ?, ?, ? , ?, ?, ?)`,
		Uuid,
		payload.Name,
		payload.Email,
		payload.PhoneNumber,
		passwordHash,
		&timeStr,
		&timeStr,
	); err != nil {
		return err
	}

	return nil
}

// GetUserByFieldName implements UserRepository.
func (r *userRepositoryImpl) GetUserByFieldName(ctx context.Context, fieldName string, value string) (*entity.User, error) {
	_query := query.SelectQuery(
		"users",
		[]string{
			"id",
			"name",
			"email",
			"password",
			"phone_number",
			"created_at",
			"updated_at",
		},
		fieldName+" = ?",
		1,
		0,
	)

	var result entity.User

	row := r.db.QueryRowX(ctx, _query, value)

	if err := row.Scan(
		&result.Id,
		&result.Name,
		&result.Email,
		&result.Password,
		&result.PhoneNumber,
		&result.CreatedAt,
		&result.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &result, nil

}

func (r userRepositoryImpl) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, opts)
}

func NewUserRepositoryImpl(db mysql.Adapter) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}
