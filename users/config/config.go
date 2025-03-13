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

// Middleware для обработки CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешённый домен
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Разрешаем куки

		// Обработка preflight-запроса (OPTIONS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func StartServer() {
	mux := http.NewServeMux()
	handler := corsMiddleware(mux)

	// Запуск сервера на определенном порту
	err := http.ListenAndServe(":" + os.Getenv("SERVER_PORT"), handler)
	if err != nil {
		log.Println("Error starting user server")
		panic("Error starting user server")
	}
}