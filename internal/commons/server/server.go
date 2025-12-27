package server

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type InternalServerSettings struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	EnablePprof  bool
}

func NewInternal(settings InternalServerSettings) *http.Server {
	mux := http.NewServeMux()

	if settings.EnablePprof {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Handler:      mux,
		ReadTimeout:  settings.ReadTimeout,
		WriteTimeout: settings.WriteTimeout,
		Addr:         settings.Address,
	}
}

func GracefulShutdown(ctx context.Context, timeout time.Duration, apps ...*http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cf := context.WithTimeout(ctx, timeout)
	defer cf()

	var wg sync.WaitGroup
	for _, app := range apps {
		wg.Go(func() {
			if err := app.Shutdown(ctx); err != nil {
				slog.Error("Error shutting down server", "error", err)
			}
		})
	}

	finished := make(chan any)
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
