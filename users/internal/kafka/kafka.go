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

		// Вызываем processTask и подготавливаем ответ
		responseData, statusCode, err := processTask(task, userRepo)
		var response models.ServiceResponse
		if err != nil {
			log.Println("Ошибка выполнения задачи:", err)
			response = models.ServiceResponse{
				Status:  "error",
				Code:    statusCode,
				Message: err.Error(),
			}
		} else {
			response = models.ServiceResponse{
				Status: "success",
				Code:   statusCode,
				Data:   responseData,
			}
		}

		// Отправляем сформированный ответ в Kafka
		sendResponseToKafka("users-response", response)
	}
}


// Поиск новых сообщений
func processTask(task map[string]interface{}, userRepo *repository.UserRepository) (interface{}, int, error) {
	action, ok := task["action"].(string)
	if !ok {
		log.Println("Ошибка: нет действия")
		return nil, 500, nil
	}

	switch action {
	case "login":
		userData := models.UserLoginRequest{
			Email:         task["email"].(string),
			InputPassword: task["password"].(string),
		}

		user, statusCode, err := userRepo.UserLogin(userData)
		if err != nil {
			log.Println("Ошибка логина:", err)
			return nil, statusCode, err
		}

		log.Println("Пользователь успешно залогинен:", user)
		return user, statusCode, nil

	case "create":
		userData := models.CreateUser{
			Email:    task["email"].(string),
			Username: task["username"].(string),
			Password: task["password"].(string),
			Gender:   task["gender"].(string),
			Location: task["location"].(string),
		}

		newUser, statusCode, err := userRepo.CreateUser(userData)
		if err != nil {
			log.Println("Ошибка создания пользователя:", err)
			return nil, statusCode, err
		}

		log.Println("Пользователь успешно создан:", newUser)
		return newUser, statusCode, nil

	default:
		log.Println("Неизвестное действие:", action)
		return nil, 400, nil
	}
}
