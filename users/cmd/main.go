package main

import (
	"log"
	"users_service/config"
	"users_service/internal/auth"
	"users_service/internal/kafka"
	"users_service/internal/repository"
)

func main() {
	// 1. Загрузка переменных окружения
	config.LoadEnv()

	// 2. Подключение к БД
	db := config.ConnectDB()

	// 3. Инициализация менеджера сессии
	sessionManager := auth.NewSessionManager(db)

	// 4. Инициализация репозитория пользователей
	userRepo := repository.NewUserRepository(db, sessionManager)

	// 5. Запуск Kafka Consumer в отдельной горутине
	go kafka.UserServiceConsumer(userRepo)

	// 6. Запуск HTTP-сервера
	log.Println("Запуск сервера пользователей...")
	config.StartServer()
}
