package usecases

import (
	"errors"
	"github.com/gmurayama/url-shortner/application"
	"github.com/gmurayama/url-shortner/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

func TestShorten(t *testing.T) {
	t.Parallel()
	input := "http://www.google.com"
	l := zap.NewNop()
	h := hash(input)

	t.Run("should return error if cannot search for shortened URL", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		dalMock := mocks.NewMockDal(ctrl)
		dalMock.EXPECT().Find(h).Return("", errors.New("error"))

		useCases := NewShortenerUseCases(dalMock, l)
		res, err := useCases.Shorten(input)

		assert.Empty(t, res)
		assert.Error(t, err)
	})

	t.Run("should return error if cannot save shortened URL", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		dalMock := mocks.NewMockDal(ctrl)
		dalMock.EXPECT().Find(h).Return("", application.ErrNotFound)

		dalMock.EXPECT().Save(h, input).Return(errors.New("error"))

		useCases := NewShortenerUseCases(dalMock, l)
		res, err := useCases.Shorten(input)

		assert.Empty(t, res)
		assert.Error(t, err)
	})

	t.Run("should shorten URL", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		dalMock := mocks.NewMockDal(ctrl)
		dalMock.EXPECT().Find(h).Return("", application.ErrNotFound)

		dalMock.EXPECT().Save(h, input).Return(nil)

		useCases := NewShortenerUseCases(dalMock, l)
		res, err := useCases.Shorten(input)

		assert.NotEmpty(t, res)
		assert.Equal(t, h, res)
		assert.NoError(t, err)
	})

	t.Run("should handle collision and shorten URL", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		dalMock := mocks.NewMockDal(ctrl)
		dalMock.EXPECT().Find(h).Return("https://notgoogle.com", nil)
		dalMock.EXPECT().Find(gomock.Any()).Return("", application.ErrNotFound)

		dalMock.EXPECT().Save(gomock.Any(), input).DoAndReturn(func(hashed, url string) error {
			assert.NotEqual(t, hashed, h)
			return nil
		})

		useCases := NewShortenerUseCases(dalMock, l)
		res, err := useCases.Shorten(input)

		assert.NotEmpty(t, res)
		assert.NoError(t, err)
	})

}
