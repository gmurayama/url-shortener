package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gmurayama/url-shortener/config"
	"github.com/gmurayama/url-shortener/internal/application"
	"github.com/gmurayama/url-shortener/internal/gateways/api/handlers"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func New(cfg *config.Config) http.Handler {
	r := gin.Default()
	r.Use(otelgin.Middleware(
		cfg.Application.Name,
		otelgin.WithGinFilter(func(c *gin.Context) bool { return c.FullPath() != "/healthz" })),
	)
	r.Use(func(c *gin.Context) {
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

	shortenUseCase := application.NewShortenUseCase()

	healthcheckHandler := handlers.NewHealthcheckHandler()
	urlHandler := handlers.NewURLHandler(shortenUseCase)

	r.GET("/healthz", healthcheckHandler.Handler)
	r.POST("/shorten", urlHandler.Shorten)

	return r
}
