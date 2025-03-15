package main

import (
	"gateway/config"
	"gateway/handler"
)

func main() {
	// 1. Запуск сервера
	app := config.StartServer()

	// 2. Инициализация роутера
	handler.Handlers(app)
}