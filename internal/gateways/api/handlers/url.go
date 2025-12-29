package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gmurayama/url-shortener/internal/application"
)

type URLHandler struct {
	shortenURLUseCase application.ShortenUseCase
	getURLUseCase     application.GetURLUseCase
}

func NewURLHandler(
	shortenUseCase application.ShortenUseCase,
	getURLUseCase application.GetURLUseCase,
) URLHandler {
	return URLHandler{
		shortenURLUseCase: shortenUseCase,
		getURLUseCase:     getURLUseCase,
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

func (h *URLHandler) GetURL(c *gin.Context) {
	shortenedURL := c.Param("shortened")
	if shortenedURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shortened URL can not be empty"})
		return
	}

	slog.Info("shortenedURL", slog.String("shortened", shortenedURL))

	url, err := h.getURLUseCase.GetURL(c.Request.Context(), shortenedURL)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrShortenedURLNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "shortened URL not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.Redirect(http.StatusPermanentRedirect, url)
}
