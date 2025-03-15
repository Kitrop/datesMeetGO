package handler

import (
	"gateway/models"
	"gateway/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func Handlers(app *fiber.App) {
	apiUserService := app.Group("/users")

	apiUserService.Post("/create", func(c *fiber.Ctx) error {
		var input models.UserCreateInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Неверный формат запроса",
			})
		}
		if err := validate.Struct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка валидации: " + err.Error(),
			})
		}

		// Создаем запрос типа KafkaRequest
		req := service.KafkaRequest{
			Action:   "create",
			Email:    input.Email,
			Username: input.Username,
			Password: input.Password,
			Gender:   input.Gender,
			Location: input.Location,
		}

		// Ждем ответ
		resp, err := service.SendRequestToUserService(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка обработки запроса: " + err.Error(),
			})
		}

		return c.Status(resp.Code).JSON(resp)
	})
	

	apiUserService.Post("/login", func(c *fiber.Ctx) error {
		var input models.UserLoginInput
		// Парсинг тела запроса в структуру
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Неверный формат запроса",
			})
		}
		// Валидация входных данных
		if err := validate.Struct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка валидации: " + err.Error(),
			})
		}

		

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Пользователь успешно авторизован",
		})
	})
}
