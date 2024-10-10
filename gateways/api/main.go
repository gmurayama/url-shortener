package main

import (
	"context"
	"fmt"
	log "log/slog"

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
	addr := fmt.Sprintf("%s:%d", cfg.Application.Address, cfg.Application.Port)
	go app.StartServer(svr, addr)

	internal := app.NewInternal(cfg)
	go app.StartServer(internal, addr)

	app.GracefulShutdown(ctx, cfg.Application.ShutdownTimeout, svr, internal)
}
