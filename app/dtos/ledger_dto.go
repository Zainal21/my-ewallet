package dtos

type TopUpTransactionDto struct {
	UserId string `json:"user_id" validate:"required"`
	RefId  string `json:"ref_id" validate:"nullable"`
	Type   string `json:"type" validate:"required"`
	Amount string `json:"amount" validate:"required"`
}

type TransferRequestDto struct {
	UserId             string `json:"user_id" validate:"required"`
	RefId              string `json:"ref_id" validate:"nullable"`
	Type               string `json:"type" validate:"required"`
	Amount             string `json:"amount" validate:"required"`
	AccountDestination string `json:"account_destination" validate:"required"`
	Bank               string `json:"bank" validate:"required"`
}

type TransactionRequestDto struct {
	UserId   string `json:"user_id" validate:"required"`
	Search   string `json:"search" validate:"nullable"`
	DateFrom string `json:"date_from" validate:"nullable"`
	DateTo   string `json:"date_to" validate:"nullable"`
	Page     int    `json:"page" validate:"required"`
}

type LedgerDto struct {
	UserID         string `json:"user_id" db:"user_id"`
	RefID          string `json:"ref_id" db:"ref_id"`
	Type           string `json:"type" db:"type"`
	CurrentDeposit int    `json:"current_deposit" db:"current_deposit"`
	ChangeDeposit  string `json:"change_deposit" db:"change_deposit"`
	FinalDeposit   int    `json:"final_deposit" db:"final_deposit"`
	Note           string `json:"note" db:"note"`
}
