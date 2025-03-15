package main

import (
	"gateway/config"
	"gateway/handler"
	"gateway/kafka"
	"gateway/service"
)

func main() {
	// 1. Запуск сервера
	app := config.StartServer()
	app.Listen(":8080")

	// 2. Инициализация роутера
	handler.Handlers(app)

	// Запуск consumer kafka
	go kafka.GatewayResponseConsumer("users-response", service.ProcessKafkaResponse)

}