package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
		log.Fatal("Error signed jwt token")
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
		fmt.Println(err.Error())
		return nil, err
	}

	if !parsedToken.Valid {
		err := errors.New("invalid token")
		return nil, err
	}

	return parsedToken.Claims, nil
}