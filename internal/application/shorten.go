package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShortenUseCase struct {
	conn *pgxpool.Pool
}

func NewShortenUseCase(conn *pgxpool.Pool) ShortenUseCase {
	return ShortenUseCase{
		conn: conn,
	}
}

func (uc *ShortenUseCase) Shorten(ctx context.Context, url string) error {
	for {
		hash := sha256.Sum256([]byte(url))
		hashString := hex.EncodeToString(hash[:])
		shortened := hashString[:7]

		_, err := uc.conn.Exec(ctx, "INSERT INTO urls(url, shortened) VALUES($1, $2)", url, shortened)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			continue
		}

		return err
	}
}
