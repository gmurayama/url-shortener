package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer(cfg Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:           cfg.Application.ReadTimeout,
		WriteTimeout:          cfg.Application.WriteTimeout,
		DisableStartupMessage: true,
		AppName:               cfg.Application.Name,
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
		AppName:               "Internal Server",
	})

	app.Use(healthcheck.New())
	app.Use(pprof.New())
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	return app
}

func StartServer(app *fiber.App, addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("error starting listener", "appName", app.Config().AppName, "address", addr)
		panic(err)
	}
	slog.Info(fmt.Sprintf("Listening on %s", addr))

	if err := app.Listener(listener); err != nil {
		slog.Error(fmt.Sprintf("error on server %s", addr), "error", err)
	}
}

func GracefulShutdown(ctx context.Context, timeout time.Duration, apps ...*fiber.App) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cf := context.WithTimeout(ctx, timeout)
	defer cf()

	var wg sync.WaitGroup
	for _, app := range apps {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := app.ShutdownWithContext(ctx); err != nil {
				slog.Error("Error shutting down server", "error", err, "appName", app.Config().AppName)
			}
		}()
	}

	finished := make(chan interface{})
	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-ctx.Done():
		slog.Error("Shutting down servers timed out", "error", ctx.Err())
	case <-finished:
		slog.Info("Stopped servers")
	}
}
