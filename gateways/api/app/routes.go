package app

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Post("/v1/create", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	app.Get("/:short", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}
