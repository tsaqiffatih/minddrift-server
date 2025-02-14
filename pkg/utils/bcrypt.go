package utils

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	upperCase, _ := regexp.MatchString("[A-Z]", password)
	lowerCase, _ := regexp.MatchString("[a-z]", password)
	digit, _ := regexp.MatchString("[0-9]", password)
	// specialChar, _ := regexp.MatchString("[!@#$%^&*(),.?\":{}|<>]", password)

	if !upperCase || !lowerCase || !digit {
		return errors.New("password must include at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}
