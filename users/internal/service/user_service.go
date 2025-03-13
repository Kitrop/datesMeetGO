package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/matthewhartstonge/argon2"
)

// Хэширование пароля
func HashPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()
	
	hash, err := argon.HashEncoded([]byte(password))
	if err != nil {
		log.Println("Error hash password")
		err := errors.New("error hash password")
		return "", err
	}

	return string(hash), nil
}

// Проверка на совпадения пароля 
func PasswordCompare(hashPassword, inputPassword string) (bool, error){
	ok, err := argon2.VerifyEncoded([]byte(inputPassword), []byte(hashPassword))

	if err != nil {
		log.Println(err)
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

type Claims struct {
	UserID uint `json:"userID"`
	Username string `json:"username"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CreateJWT(userID uint, username string, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Username: username,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(os.Getenv("JWT_KEY"))

	if err != nil {
		log.Println("Error signed jwt token")
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(jwtToken string) (jwt.Claims, error){
	// Парсинг и проверка токена на валидность
	parsedToken, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется HMAC метод
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неправильный метод подписания: %v", token.Header["alg"])
		}
		return os.Getenv("JWT_KEY"), nil
	})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !parsedToken.Valid {
		err := errors.New("invalid token")
		return nil, err
	}

	return parsedToken.Claims, nil
}