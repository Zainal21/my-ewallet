package dtos

type UserSignInRequestDto struct {
	Email    string `json:"email" validate:"required|email"`
	Password string `json:"password" validate:"required"`
}

type UserSignInResponseDto struct {
	Id          string  `json:"id" db:"id"`
	Email       string  `json:"email" db:"email"`
	Name        string  `json:"name" db:"name"`
	PhoneNumber *string `json:"phone_number" db:"phone_number"`
	Token       string  `json:"token" db:"token"`
	CreatedAt   *string `json:"created_at" db:"created_at"`
	UpdatedAt   *string `json:"updated_at" db:"updated_at"`
}
