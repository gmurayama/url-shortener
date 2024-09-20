package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Config struct {
	Application struct {
		Address string `env:"ADDR" envDefault:"localhost"`
		Port    int    `env:"PORT" envDefault:"8080"`
	} `envPrefix:"APP_"`
}

func main() {
	cfg, err := env.ParseAsWithOptions[Config](env.Options{
		Prefix: "SH_",
	})
	if err != nil {
		log.Panic("could not load config", err)
	}

	app := fiber.New()
	app.Use(recover.New())

	app.Post("/v1/create", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	app.Get("/:short", func(c *fiber.Ctx) error {
		return nil
	})

	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Application.Address, cfg.Application.Port)
		log.Println("Listening on", addr)
		if err := app.Listen(addr); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	log.Println("Running cleanup tasks...")
}
