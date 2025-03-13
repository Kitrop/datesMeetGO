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

	// Создание сессии и токена
	token, err := CreateSession(db, newUser.ID, newUser.Username, newUser.Email)
	if err != nil {
		return newUserData{}, err
	}

	return newUserData{UserID: newUser.ID, Email: newUser.Email, Username: newUser.Username, token: token}, nil
}

// Создание сессии и токена для пользователя
func CreateSession(db *gorm.DB, userID uint, username string, email string) (string, error) {
	// Cоздание токена
	token, err := service.CreateJWT(userID, username, email)
	if err != nil {
		return "", err
	}

	newSession := models.UserSession{
		UserID: userID,
		Token: token,
	}

	// Создание сессии в БД
	result := db.Create(&newSession)
	if result.Error != nil {
		log.Fatal("Ошибка при создании токена:", result.Error)
		return "", result.Error
	}

	return token, nil
}

// Получение пользователя из БД
func getUserFromDB(db *gorm.DB, userID uint) (models.User, error){
	var user models.User
	findUser := db.First(&user, userID) // Поиск пользователя по id

	// Если пользователь не найден
	if findUser.Error != nil {
		err := errors.New("user not found")
		return models.User{}, err
	}

	return user, nil
}

// Функция для логина пользователя
func UserLogin(db *gorm.DB, userID uint, inputPassword string) (models.UserLoginResponse, error) {
	userData, err := getUserFromDB(db, userID) // Получаем данные 
	if err != nil {
		return models.UserLoginResponse{}, err
	}

	// Проверка на правильность пароля
	isLogin, err := service.PasswordCompare(userData.PasswordHash, inputPassword)
	if err != nil {
		return models.UserLoginResponse{}, err
	}

	if !isLogin {
		err := errors.New("incorrect password")
		return models.UserLoginResponse{}, err
	}

	token, err := CreateSession(db, userData.ID, userData.Username, userData.Email)
	if err != nil {
		return models.UserLoginResponse{}, err
	}

	return models.UserLoginResponse{ Username: userData.Username, Email: userData.Email, Token: token}, nil
}