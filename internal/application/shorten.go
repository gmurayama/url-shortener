package application

import (
	"crypto/sha256"
	"encoding/hex"
)

type ShortenUseCase struct{}

func NewShortenUseCase() ShortenUseCase {
	return ShortenUseCase{}
}

func (uc *ShortenUseCase) Shorten(url string) string {
	hash := sha256.Sum256([]byte(url))
	hashString := hex.EncodeToString(hash[:])

	return hashString[:7]
}
