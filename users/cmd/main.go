package main

import (
	"users_service/config"
)

func main() {
	config.LoadEnv() // Загрузка env
	config.ConnectDB() // Подключение к БД
	config.StartServer() // Запуск сервера
}