package main

import (
	"context"
	"fmt"
	log "log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/gmurayama/url-shortner/gateways/api/app"
	"github.com/gmurayama/url-shortner/infrastructure/tracing"
)

func main() {
	ctx := context.Background()

	cfg, err := env.ParseAsWithOptions[app.Config](env.Options{
		Prefix: "SH_",
	})
	if err != nil {
		log.Error("could not load config", "error", err)
		panic(err)
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
		log.Error("could not set tracing", "error", err)
		panic(err)
	}
	defer s(ctx)

	svr := app.NewServer(cfg)
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Application.Address, cfg.Application.Port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Error("error starting listener", "address", addr)
			panic(err)
		}
		log.Info(fmt.Sprintf("Listening on %s", addr))

		if err := svr.Listener(listener); err != nil {
			log.Error("error on server", "error", err)
		}
	}()

	internal := app.NewInternal(cfg)
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Internal.Address, cfg.Internal.Port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Error("error starting listener", "address", addr)
			panic(err)
		}
		log.Info(fmt.Sprintf("Listening on %s", addr))

		if err := internal.Listener(listener); err != nil {
			log.Error("error on internal server", "error", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Gracefully shutting down...")
	_ = svr.Shutdown()

	log.Info("Running cleanup tasks...")
}
