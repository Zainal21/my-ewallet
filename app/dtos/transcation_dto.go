package dtos

type TransactionDto struct {
	OrderID     string  `json:"order_id" db:"order_id"`
	UserID      string  `json:"user_id" db:"user_id"`
	RefID       string  `json:"ref_id" db:"ref_id"`
	Type        string  `json:"type" db:"type"`
	GrossAmount string  `json:"gross_amount" db:"gross_amount"`
	Piece       string  `json:"piece" db:"piece"`
	Amount      string  `json:"amount" db:"amount"`
	Note        string  `json:"note" db:"note"`
	Status      string  `json:"status" db:"status"`
	CreatedAt   *string `json:"created_at" db:"created_at"`
	UpdatedAt   *string `json:"updated_at" db:"updated_at"`
}
