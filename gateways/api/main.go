package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/gmurayama/url-shortner/infrastructure/tracing"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Config struct {
	Application struct {
		Name    string `env:"NAME" envDefault:"URL Shortener"`
		Address string `env:"ADDR" envDefault:"localhost"`
		Port    int    `env:"PORT" envDefault:"8080"`
	} `envPrefix:"APP_"`
	Tracing struct {
		Host               string        `env:"HOST" envDefault:"localhost"`
		Port               int           `env:"PORT" envDefault:"4317"`
		Enabled            bool          `env:"ENABLED" envDefault:"true"`
		BatchScheduleDelay time.Duration `env:"BATCH_SCHEDULE_DELAY" envDefault:"5s"`
		SamplingRatio      float64       `env:"SAMPLING_RATIO" envDefault:"1.0"`
		MaxExportBatchSize int           `env:"MAX_EXPORT_BATCH_SIZE" envDefault:"256"`
		KeepAliveTime      time.Duration `env:"KEEP_ALIVE_TIME" envDefault:"20s"`
		KeepAliveTimeout   time.Duration `env:"KEEP_ALIVE_TIMEOUT" envDefault:"5s"`
	} `envPrefix:"TRACING_"`
}

func main() {
	ctx := context.Background()

	cfg, err := env.ParseAsWithOptions[Config](env.Options{
		Prefix: "SH_",
	})
	if err != nil {
		log.Panic("could not load config", err)
	}

	s, err := tracing.Configure(ctx, tracing.Settings{
		ServiceName:        cfg.Application.Name,
		Host:               cfg.Tracing.Host,
		Port:               cfg.Tracing.Port,
		Enabled:            cfg.Tracing.Enabled,
		BatchScheduleDelay: cfg.Tracing.BatchScheduleDelay,
		SamplingRatio:      cfg.Tracing.SamplingRatio,
		MaxExportBatchSize: cfg.Tracing.MaxExportBatchSize,
		KeepAliveTime:      cfg.Tracing.KeepAliveTime,
		KeepAliveTimeout:   cfg.Tracing.KeepAliveTimeout,
	})
	if err != nil {
		log.Panic("could not set tracing", err)
	}
	defer s(ctx)

	app := fiber.New()

	app.Use(otelfiber.Middleware())
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
