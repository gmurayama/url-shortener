package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gmurayama/url-shortener/internal/application"
)

type URLHandler struct {
	shortenURLUseCase application.ShortenUseCase
}

func NewURLHandler(shortenUseCase application.ShortenUseCase) URLHandler {
	return URLHandler{
		shortenURLUseCase: shortenUseCase,
	}
}

type ShortenRequest struct {
	URL string `binding:"url"`
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shorten, err := h.shortenURLUseCase.Shorten(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shorten": shorten})
}
