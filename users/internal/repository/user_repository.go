package repository

import (
	"errors"
	"log"
	"users_service/internal/models"
	"users_service/internal/service"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type newUserData struct {
	UserID uint
	Username string
	Email string
	token string
}

// Создание нового пользователя
func CreateUser(db *gorm.DB, userData models.CreateUser) (newUserData, error) {
	// Валидация входных данных
	validate := validator.New()
	err := validate.Struct(userData)
	if err != nil {
		return newUserData{}, err
	}
	
	// Проверка пароля на сложность
	isStrong := service.IsStrongPassword(userData.Password)
	if !isStrong {
		err := errors.New("not strong password")
		return newUserData{}, err
	}

	// Хэширование пароля
	hashPassword, err := service.HashPassword(userData.Password)
	if err != nil {
		return newUserData{}, err
	}

	newUser := models.User{
		Email: userData.Email,
		PasswordHash: hashPassword,
		Username: userData.Username,
		Gender: userData.Gender,
		Location: userData.Location,
		BirthDate: userData.BirthDate,
	}

	// Создание пользователя в системе
	result := db.Create(&newUser)
	if result.Error != nil {
		log.Fatal("Ошибка при создании пользователя:", result.Error)
		return newUserData{}, result.Error
	}

	// Cоздание токена
	token, err := service.CreateJWT(newUser.ID, newUser.Username, newUser.Email)
	if err != nil {
		return newUserData{}, err
	}

	newSession := models.UserSession{
		UserID: newUser.ID,
		Token: token,
	}
	
	// Создание сессии для пользователя
	resultNewSession := db.Create(&newSession)
	if resultNewSession.Error != nil {
		log.Fatal("Ошибка при создании токена:", result.Error)
		return newUserData{}, result.Error
	}

	return newUserData{UserID: newUser.ID, Email: newUser.Email, Username: newUser.Username, token: token}, nil
}

func UserLogin(password string) {
	
}