package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// Отправляет запрос в указанный топик
func SendRequestToKafka(requestTopic string, request interface{}) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   requestTopic,
	})
	defer writer.Close()

	data, err := json.Marshal(request)
	if err != nil {
		log.Println("Ошибка кодирования JSON:", err)
		return err
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{Value: data})
	if err != nil {
		log.Println("Ошибка отправки сообщения в Kafka:", err)
		return err
	}
	return nil
}

// Читает сообщения из топика ответов и вызывает обработчик для каждого полученного сообщения
func GatewayResponseConsumer(responseTopic string, handler func(message []byte)) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   responseTopic,
		GroupID: "gateway-response-group",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Ошибка чтения сообщения:", err)
			continue
		}
		// Вызов обработчика для дальнейшей обработки полученного сообщения
		handler(msg.Value)
	}
}