package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/gmurayama/url-shortener/config"
	"github.com/gmurayama/url-shortener/internal/commons/server"
	"github.com/gmurayama/url-shortener/internal/gateways/api"
	"github.com/gmurayama/url-shortener/internal/infrastructure/tracing"
)

func main() {
	if err := run(); err != nil {
		slog.Error("exec error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return nil
		}

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

	apiImpl, err := api.New(ctx, &cfg)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:         cfg.Server.Address,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		Handler:      apiImpl,
	}
	internalSrv := server.NewInternal(server.InternalServerSettings{
		Address:      cfg.InternalServer.Address,
		WriteTimeout: cfg.InternalServer.WriteTimeout,
		ReadTimeout:  cfg.InternalServer.ReadTimeout,
		EnablePprof:  cfg.InternalServer.EnablePprof,
	})

	api.RegisterMetrics()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error on server", "error", err)
		}
	}()
	go func() {
		if err := internalSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error on internal server", "error", err)
		}
	}()
	server.GracefulShutdown(ctx, cfg.Server.ShutdownTimeout, srv, internalSrv)

	return nil
}
