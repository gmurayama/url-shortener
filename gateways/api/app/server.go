package app

import (
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

func NewServer(cfg Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:           cfg.Application.ReadTimeout,
		WriteTimeout:          cfg.Application.WriteTimeout,
		DisableStartupMessage: true,
	})

	app.Use(otelfiber.Middleware())
	app.Use(recover.New())
	app.Use(timeout.NewWithContext(func(c *fiber.Ctx) error { return c.Next() }, cfg.Application.Timeout))

	RegisterRoutes(app)

	return app
}

func NewInternal(cfg Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:           cfg.Internal.ReadTimeout,
		WriteTimeout:          cfg.Internal.WriteTimeout,
		DisableStartupMessage: true,
	})

	app.Use(healthcheck.New())
	app.Use(pprof.New())

	return app
}
