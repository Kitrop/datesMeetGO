package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"users_service/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDSN() string {
	// Генерация ссылки подключения к БД
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	
	return dsn
}

func ConnectDB() *gorm.DB {
	// Создание подключения
	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	// Создаем расширение uuid-ossp, если его нет
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// Миграция схем
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.UserSession{})

	return db
}

func LoadEnv() {
	// Загрузка env
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
		panic("Error loading .env file")
	}
}


func StartServer() {
	// Запуск сервера на определенном порту
	err := http.ListenAndServe(":" + os.Getenv("SERVER_PORT"), nil)
	if err != nil {
		log.Println("Error starting user server")
		panic("Error starting user server")
	}
}