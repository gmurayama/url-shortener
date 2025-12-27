package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gmurayama/webservice-template-golang/config"
	"github.com/gmurayama/webservice-template-golang/internal/gateways/api/handlers"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func New(cfg *config.Config) http.Handler {
	healthcheckHandler := handlers.NewHealthcheckHandler()

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

	r.GET("/healthz", healthcheckHandler.Handler)

	return r
}
