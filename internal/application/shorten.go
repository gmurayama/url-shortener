package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

type ShortenUseCase struct {
	conn *pgxpool.Pool

	metrics ShortenUseCaseMetrics
}

type ShortenUseCaseMetrics struct {
	shortenerAttempts *prometheus.CounterVec
}

func NewShortenUseCase(conn *pgxpool.Pool) ShortenUseCase {
	return ShortenUseCase{
		conn: conn,

		metrics: ShortenUseCaseMetrics{
			shortenerAttempts: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "shortener_attempts",
				Help: "How many attempts of registration of a short URL it took",
			},
				[]string{"attempt"},
			),
		},
	}
}

func (uc *ShortenUseCase) RegisterMetrics() {
	prometheus.MustRegister(uc.metrics.shortenerAttempts)
}

// TODO: how to handle if the URL is already registered? return the shortened URL or emit an error?
func (uc *ShortenUseCase) Shorten(ctx context.Context, url string) (string, error) {
	attempts := new(int)

	defer func() {
		uc.metrics.shortenerAttempts.With(prometheus.Labels{"attempt": fmt.Sprintf("%d", *attempts)}).Inc()
	}()

	urlToHash := url
	for *attempts = 1; *attempts < 5; *attempts++ {
		if *attempts > 1 {
			asciiChar := rand.IntN(126-33) + 33
			urlToHash = fmt.Sprintf("%s%c", urlToHash, rune(asciiChar))
		}

		h := sha256.New()
		h.Write([]byte(urlToHash))
		hash := h.Sum(nil)
		hashString := hex.EncodeToString(hash[:])
		shortened := hashString[:7]

		slog.Debug("generated urlHash",
			slog.String("urlHash", urlToHash),
			slog.String("hash", shortened),
			slog.Int("attempt", *attempts),
		)

		_, err := uc.conn.Exec(ctx, "INSERT INTO urls(url, shortened) VALUES($1, $2)", url, shortened)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			continue
		}

		return shortened, err
	}

	return "", fmt.Errorf("url shortening failed, max attempts reached")
}
