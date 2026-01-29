package v1

import (
	"github.com/andreyxaxa/URL-Shortener/internal/usecase"
	"github.com/andreyxaxa/URL-Shortener/pkg/logger"
)

type V1 struct {
	lk usecase.Link
	l  logger.Interface

	baseURL string
}
