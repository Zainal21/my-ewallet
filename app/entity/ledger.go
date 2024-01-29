package entity

type Ledger struct {
	Id             string  `json:"id" db:"id"`
	UserID         string  `json:"user_id" db:"user_id"`
	RefID          string  `json:"ref_id" db:"ref_id"`
	Type           string  `json:"type" db:"type"`
	CurrentDeposit string  `json:"current_deposit" db:"current_deposit"`
	ChangeDeposit  string  `json:"change_deposit" db:"change_deposit"`
	FinalDeposit   string  `json:"final_deposit" db:"final_deposit"`
	Note           string  `json:"note" db:"note"`
	CreatedAt      *string `json:"created_at" db:"created_at"`
	UpdatedAt      *string `json:"updated_at" db:"updated_at"`
}
