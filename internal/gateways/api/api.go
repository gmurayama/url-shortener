package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gmurayama/url-shortener/config"
	"github.com/gmurayama/url-shortener/internal/application"
	"github.com/gmurayama/url-shortener/internal/gateways/api/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func New(ctx context.Context, cfg *config.Config) (http.Handler, error) {
	r := gin.Default()

	RegisterMiddlewares(r, cfg)

	pgConn, err := pgxpool.New(ctx, cfg.Database.ConnString)
	if err != nil {
		return nil, err
	}
	shortenUseCase := application.NewShortenUseCase(pgConn)
	shortenUseCase.RegisterMetrics()

	healthcheckHandler := handlers.NewHealthcheckHandler()
	urlHandler := handlers.NewURLHandler(shortenUseCase)

	r.GET("/healthz", healthcheckHandler.Handler)
	r.POST("/shorten", urlHandler.Shorten)

	return r, nil
}

func RegisterMiddlewares(router *gin.Engine, cfg *config.Config) {
	router.Use(otelgin.Middleware(
		cfg.Application.Name,
		otelgin.WithGinFilter(func(c *gin.Context) bool { return c.FullPath() != "/healthz" })),
	)
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		requestDuration.With(prometheus.Labels{
			"method": c.Request.Method,
			"route":  c.FullPath(),
			"status": fmt.Sprint(c.Writer.Status()),
		}).
			Observe(float64(duration.Milliseconds()))
	})
}
