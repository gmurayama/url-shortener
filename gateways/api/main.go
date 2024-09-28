package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/gmurayama/url-shortner/gateways/api/server"
	"github.com/gmurayama/url-shortner/infrastructure/tracing"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

type Config struct {
	Application struct {
		Name         string        `env:"NAME" envDefault:"URL Shortener"`
		Address      string        `env:"ADDR" envDefault:"localhost"`
		Port         int           `env:"PORT" envDefault:"8080"`
		ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"2s"`
		WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"2s"`
		Timeout      time.Duration `env:"TIMEOUT" envDefault:"2s"`
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

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Application.ReadTimeout,
		WriteTimeout: cfg.Application.WriteTimeout,
	})

	app.Use(otelfiber.Middleware())
	app.Use(recover.New())
	app.Use(timeout.NewWithContext(func(c *fiber.Ctx) error { return c.Next() }, cfg.Application.Timeout))
	server.RegisterRoutes(app)

	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Application.Address, cfg.Application.Port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Panic(err)
		}
		log.Println("Listening on", addr)

		if err := app.Listener(listener); err != nil {
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
