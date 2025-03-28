package utils

import (
	"fmt"
	"net/mail"
	"unicode"
)

func CheckEmailAndPassword(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("email is not valid")
	}

	return nil
}

func IsValidPassword(password string) bool {
	var (
		hasMinLen    = false
		hasUppercase = false
		hasLowercase = false
		hasNumber    = false
		hasSpecial   = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUppercase && hasLowercase && hasNumber && hasSpecial
}
