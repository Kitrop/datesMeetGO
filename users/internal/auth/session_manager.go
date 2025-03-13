package auth

import (
	"log"
	"users_service/internal/models"
	"users_service/internal/service"

	"gorm.io/gorm"
)

// SessionManager отвечает за создание JWT-токена и сохранение сессии в БД
type SessionManager struct {
	db *gorm.DB
}

// NewSessionManager создаёт новый экземпляр SessionManager
func NewSessionManager(db *gorm.DB) *SessionManager {
	return &SessionManager{db: db}
}

// CreateSession создаёт JWT-токен и сохраняет сессию в базе данных
func (sm *SessionManager) CreateSession(userID uint, username, email string) (string, error) {
	// Создание JWT токена с использованием логики из сервиса
	token, err := service.CreateJWT(userID, username, email)
	if err != nil {
		return "", err
	}

	newSession := models.UserSession{
		UserID: userID,
		Token:  token,
	}

	// Сохранение сессии в БД
	if err := sm.db.Create(&newSession).Error; err != nil {
		log.Println("Ошибка при создании сессии:", err)
		return "", err
	}

	return token, nil
}
