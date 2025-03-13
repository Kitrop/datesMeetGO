package repository

import (
	"errors"
	"log"
	"users_service/internal/models"
	"users_service/internal/auth"
	"users_service/internal/service"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserRepository struct {
	db             *gorm.DB
	validate       *validator.Validate
	sessionManager *auth.SessionManager
}

// Конструктор, принимающий sessionManager как зависимость
func NewUserRepository(db *gorm.DB, sessionManager *auth.SessionManager) *UserRepository {
	return &UserRepository{
		db:             db,
		validate:       validator.New(),
		sessionManager: sessionManager,
	}
}

// Создание нового пользователя
func (r *UserRepository) CreateUser(userData models.CreateUser) (models.UserCreateResponse, error) {
	// Валидация входных данных
	if err := r.validate.Struct(userData); err != nil {
		return models.UserCreateResponse{}, err
	}

	// Проверка пароля на сложность
	if !service.IsStrongPassword(userData.Password) {
		return models.UserCreateResponse{}, errors.New("password is not strong enough")
	}

	// Хэширование пароля
	hashPassword, err := service.HashPassword(userData.Password)
	if err != nil {
		return models.UserCreateResponse{}, err
	}

	newUser := models.User{
		Email:        userData.Email,
		PasswordHash: hashPassword,
		Username:     userData.Username,
		Gender:       userData.Gender,
		Location:     userData.Location,
		BirthDate:    userData.BirthDate,
	}

	// Создание пользователя в системе
	if err := r.db.Create(&newUser).Error; err != nil {
		// Проверка дублирования
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return models.UserCreateResponse{}, errors.New("user with this email already exists")
		}
		log.Println("Ошибка при создании пользователя:", err)
		return models.UserCreateResponse{}, err
	}

	// Создание сессии и токена через sessionManager
	token, err := r.sessionManager.CreateSession(newUser.ID, newUser.Username, newUser.Email)
	if err != nil {
		return models.UserCreateResponse{}, err
	}

	return models.UserCreateResponse{
		Username: newUser.Username,
		Email:    newUser.Email,
		Token:    token,
	}, nil
}

// Получение пользователя из БД
func (r *UserRepository) GetUserFromDB(userID uint) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return user, nil
}

// Логин пользователя
func (r *UserRepository) UserLogin(userDataInput models.UserLoginRequest) (models.UserLoginResponse, error) {
	// Валидация входных данных
	if err := r.validate.Struct(userDataInput); err != nil {
		return models.UserLoginResponse{}, err
	}

	// Получаем данные пользователя из БД
	userData, err := r.GetUserFromDB(userDataInput.UserID)
	if err != nil {
		return models.UserLoginResponse{}, err
	}

	// Проверка пароля
	if isValid, err := service.PasswordCompare(userData.PasswordHash, userDataInput.InputPassword); err != nil || !isValid {
		return models.UserLoginResponse{}, errors.New("incorrect password")
	}

	// Создание сессии и токена через sessionManager
	token, err := r.sessionManager.CreateSession(userData.ID, userData.Username, userData.Email)
	if err != nil {
		return models.UserLoginResponse{}, err
	}

	return models.UserLoginResponse{
		Username: userData.Username,
		Email:    userData.Email,
		Token:    token,
	}, nil
}