package helpers // Replace with the actual package name

import (
	"testing"
)

func TestVerifyPassword(t *testing.T) {
	hashedPassword, err := HashPassword("testPassword")
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	tests := []struct {
		name           string
		passwordHash   string
		password       string
		expectedResult bool
	}{
		{"CorrectPassword", hashedPassword, "testPassword", true},
		{"IncorrectPassword", hashedPassword, "wrongPassword", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VerifyPassword(tt.passwordHash, tt.password)
			if result != tt.expectedResult {
				t.Errorf("Expected result: %v, but got: %v", tt.expectedResult, result)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "testPassword"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	if hashedPassword == password {
		t.Error("Hashed password should not be equal to the original password")
	}

	result := VerifyPassword(hashedPassword, password)
	if !result {
		t.Error("Hashed password verification failed")
	}
}
