package application

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrShortenedURLNotFound = fmt.Errorf("shortened URL not found")

type GetURLUseCase struct {
	conn *pgxpool.Pool
}

func NewGetURLUseCase(conn *pgxpool.Pool) GetURLUseCase {
	return GetURLUseCase{conn}
}

func (uc *GetURLUseCase) GetURL(ctx context.Context, shortenedURL string) (string, error) {
	row := uc.conn.QueryRow(ctx, "SELECT url FROM urls WHERE shortened = $1", shortenedURL)
	var url string
	err := row.Scan(&url)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", ErrShortenedURLNotFound
		}
		return "", err
	}

	return url, nil
}
