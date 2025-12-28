package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gmurayama/url-shortener/internal/application"
)

type URLHandler struct {
	shortenUrlUseCase application.ShortenUseCase
}

func NewURLHandler(shortenUseCase application.ShortenUseCase) URLHandler {
	return URLHandler{
		shortenUrlUseCase: shortenUseCase,
	}
}

type ShortenRequest struct {
	Url string `binding:"url"`
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shorten := h.shortenUrlUseCase.Shorten(req.Url)

	c.JSON(http.StatusOK, gin.H{"shorten": shorten})
}
