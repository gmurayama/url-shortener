package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/gmurayama/url-shortner/gateways/api/server"
	"github.com/gmurayama/url-shortner/infrastructure/tracing"
)

func main() {
	var code int
	if err := run(); err != nil {
		code = 1
	}

	os.Exit(code)
}

func run() error {
	ctx := context.Background()

	cfg, err := env.ParseAsWithOptions[server.Config](env.Options{
		Prefix: "SH_",
	})
	if err != nil {
		slog.Error("could not load config", "error", err)
		return err
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
		slog.Error("could not set tracing", "error", err)
		return err
	}
	defer s(ctx)

	svr := server.New(cfg)
	addr := fmt.Sprintf("%s:%d", cfg.Application.Address, cfg.Application.Port)
	go server.Start(svr, addr)

	addr = fmt.Sprintf("%s:%d", cfg.Internal.Address, cfg.Internal.Port)
	internal := server.NewInternal(cfg)
	go server.Start(internal, addr)

	server.GracefulShutdown(ctx, cfg.Application.ShutdownTimeout, svr, internal)

	return nil
}
