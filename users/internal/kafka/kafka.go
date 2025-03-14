package kafka

import (
	"context"
	"encoding/json"
	"log"
	"users_service/internal/models"
	"users_service/internal/repository"

	"github.com/segmentio/kafka-go"
)

// Kafka producer для отправки ответов
func sendResponseToKafka(responseTopic string, response interface{}) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   responseTopic,
	})
	defer writer.Close()

	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Ошибка кодирования JSON:", err)
		return
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{Value: data})
	if err != nil {
		log.Println("Ошибка отправки сообщения в Kafka:", err)
	}
}

func UserServiceConsumer(userRepo *repository.UserRepository) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "users",
		GroupID: "users-service-group",
	})

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Ошибка чтения сообщения:", err)
			continue
		}

		var task map[string]interface{}
		if err := json.Unmarshal(msg.Value, &task); err != nil {
			log.Println("Ошибка разбора JSON:", err)
			continue
		}

		// Вызываем processTask и отправляем ответ в Kafka
		response, err := processTask(task, userRepo)
		if err != nil {
			log.Println("Ошибка выполнения задачи:", err)
			continue
		}

		// Определяем responseTopic (например, users-response)
		sendResponseToKafka("users-response", response)
	}
}

// processTask теперь возвращает результат (пользователя) и ошибку
func processTask(task map[string]interface{}, userRepo *repository.UserRepository) (interface{}, error) {
	action, ok := task["action"].(string)
	if !ok {
		log.Println("Ошибка: нет действия")
		return nil, nil
	}

	switch action {
	case "login":
		userData := models.UserLoginRequest{
			UserID:        uint(task["userID"].(float64)), // Приведение к uint
			InputPassword: task["password"].(string),
		}

		user, err := userRepo.UserLogin(userData)
		if err != nil {
			log.Println("Ошибка логина:", err)
			return nil, err
		}

		log.Println("Пользователь успешно залогинен:", user)
		return user, nil // ✅ Возвращаем пользователя и nil-ошибку

	case "create":
		userData := models.CreateUser{
			Email:    task["email"].(string),
			Username: task["username"].(string),
			Password: task["password"].(string),
			Gender:   task["gender"].(string),
			Location: task["location"].(string),
		}

		newUser, err := userRepo.CreateUser(userData)
		if err != nil {
			log.Println("Ошибка создания пользователя:", err)
			return nil, err
		}

		log.Println("Пользователь успешно создан:", newUser)
		return newUser, nil // ✅ Возвращаем нового пользователя и nil-ошибку

	default:
		log.Println("Неизвестное действие:", action)
		return nil, nil
	}
}
