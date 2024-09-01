package usecases

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"time"

	"github.com/gmurayama/url-shortner/application"
)

const urlLength = 7

type useCasesImpl struct {
	dal    application.Dal
	logger *zap.Logger
}

var _ application.UseCases = (*useCasesImpl)(nil)

func NewShortenerUseCases(
	dal application.Dal,
	logger *zap.Logger,
) application.UseCases {
	return &useCasesImpl{
		dal:    dal,
		logger: logger.With(zap.String("source", "useCasesImpl")),
	}
}

func (u *useCasesImpl) Shorten(url string) (string, error) {
	l := u.logger.With(
		zap.String("method", "Shorten"),
		zap.String("url", url),
	)

	// TODO: remove http:// and https:// from url

	input := url
	h := ""
	for {
		h = hash(input)
		u, err := u.dal.Find(h)
		if err != nil {
			if errors.Is(err, application.ErrNotFound) {
				l.With(zap.String("short", input)).
					Debug(fmt.Sprintf("no match found for short url, will use hash %s", h))

				break
			}

			return "", err
		}
		if u == url {
			l.With(zap.String("short", h)).
				Debug("specified url already shortened in database")

			return u, nil
		}

		input = randomChar() + url
	}

	err := u.dal.Save(h, url)
	if err != nil {
		l.Error("could not save short url to database", zap.Error(err))
		return "", err
	}

	return h, nil
}

func (u *useCasesImpl) Expand(shortened string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func hash(input string) string {
	sha := sha256.Sum256([]byte(input))
	hashString := hex.EncodeToString(sha[:])

	return hashString[:urlLength]
}

func randomChar() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(charset))
	return string(charset[randomIndex])
}
