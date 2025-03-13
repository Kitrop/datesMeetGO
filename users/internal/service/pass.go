package service

import (
	"errors"
	"log"
	"regexp"
	"unicode"

	"github.com/matthewhartstonge/argon2"
)

// Хэширование пароля
func HashPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()
	
	hash, err := argon.HashEncoded([]byte(password))
	if err != nil {
		log.Fatal("Error hash password")
		err := errors.New("Error hash password")
		return "", err
	}

	return string(hash), nil
}

// Проверка на совпадения пароля 
func PasswordCompare(hashPassword, inputPassword string) (bool, error){
	ok, err := argon2.VerifyEncoded([]byte(inputPassword), []byte(hashPassword))

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	if !ok {
		return false, nil
	} else {
		return true, nil
	}
}

// Проверка на сложность пароля
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper, hasLower, hasDigit, hasSpecial := false, false, false, false
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(string(ch)):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}