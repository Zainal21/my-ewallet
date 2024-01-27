package entity

type User struct {
	Id          string  `json:"id" db:"id"`
	Email       string  `json:"email" db:"email"`
	Name        string  `json:"name" db:"name"`
	Password    string  `json:"password" db:"password"`
	PhoneNumber *string `json:"phone_number" db:"phone_number"`
	CreatedAt   *string `json:"created_at" db:"created_at"`
	UpdatedAt   *string `json:"updated_at" db:"updated_at"`
}
