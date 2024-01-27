package helpers

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(password_hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
