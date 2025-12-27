package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthcheckHandler struct{}

func NewHealthcheckHandler() HealthcheckHandler {
	return HealthcheckHandler{}
}

func (h *HealthcheckHandler) Handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
