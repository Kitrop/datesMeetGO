package service

import (
	"encoding/json"
	"errors"
	"gateway/kafka"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type KafkaRequest struct {
	CorrelationID string `json:"correlation_id"`
	Action        string `json:"action"`
	Email         string `json:"email,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Gender        string `json:"gender,omitempty"`
	Location      string `json:"location,omitempty"`
}

type ServiceResponse struct {
	Status        string      `json:"status"`
	Code          int         `json:"code"`
	Message       string      `json:"message,omitempty"`
	Data          interface{} `json:"data,omitempty"`
	CorrelationID string      `json:"correlation_id,omitempty"`
}


var responseChannels sync.Map // map[string]chan ServiceResponse

// SendRequestToUserService отправляет запрос в Kafka и ждёт ответа.
func SendRequestToUserService(request KafkaRequest) (ServiceResponse, error) {
	// Если correlation_id не задан, генерируем его.
	if request.CorrelationID == "" {
		request.CorrelationID = uuid.New().String()
	}

	// Создаем канал для ожидания ответа.
	respChan := make(chan ServiceResponse)
	responseChannels.Store(request.CorrelationID, respChan)
	defer responseChannels.Delete(request.CorrelationID)

	// Отправляем запрос в Kafka.
	err := kafka.SendRequestToKafka("users", request)
	if err != nil {
		return ServiceResponse{}, err
	}

	// Ожидаем ответ с таймаутом.
	select {
	case resp := <-respChan:
		return resp, nil
	case <-time.After(10 * time.Second):
		return ServiceResponse{}, errors.New("таймаут ожидания ответа от микросервиса")
	}
}

// ProcessKafkaResponse вызывается из Kafka consumer при получении ответа.
func ProcessKafkaResponse(message []byte) {
	var resp ServiceResponse
	if err := json.Unmarshal(message, &resp); err != nil {
		log.Println("Ошибка разбора ответа из Kafka:", err)
		return
	}

	// Ищем канал по correlation_id.
	if ch, ok := responseChannels.Load(resp.CorrelationID); ok {
		responseChan := ch.(chan ServiceResponse)
		responseChan <- resp
	} else {
		log.Println("Не найден канал для correlation_id:", resp.CorrelationID)
	}
}
