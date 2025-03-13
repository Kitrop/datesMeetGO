package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uint         `gorm:"primaryKey"`
	Email        string       `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string       `gorm:"not null"`
	Username     string       `gorm:"size:100;not null"`
	BirthDate    time.Time    `gorm:"type:date;not null"`
	Gender       string       `gorm:"size:10;not null"`
	Location     string       `gorm:"size:255"`
	CreatedAt    time.Time    `gorm:"autoCreateTime"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime"`
	Sessions []UserSession  `gorm:"foreignKey:UserID"`
}

type UserSession struct {
	SessionID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"not null"`
}

type UserCreateResponse struct {
	ID           uint
	Email        string
	Username     string
	Token 			 string
}

type CreateUser struct {
	Username  string    `validate:"required"`
	Email     string    `validate:"required,email"`
	BirthDate time.Time `validate:"required"`
	Gender    string    `validate:"required,oneof=male female"`
	Location  string    `validate:"required"`
	Password  string    `validate:"required"`
}

type UserLoginResponse struct {
	Username  string
	Email     string
	Token 		string
}

type UserLoginRequest struct {
	UserID 				uint 	 `validate:"required"`
	InputPassword string `validate:"required"`
}