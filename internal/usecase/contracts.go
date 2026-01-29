package usecase

import (
	"context"

	"github.com/andreyxaxa/URL-Shortener/internal/entity"
)

type (
	Link interface {
		CreateShortURL(ctx context.Context, originalURL string, customAlias string) (string, error)
		GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error)
		TrackClick(ctx context.Context, shortCode, IP, userAgent string) error
		ExistsByShortCode(ctx context.Context, shortCode string) error
		GetAnalytics(ctx context.Context, shortCode string) (entity.Analytics, error)
		GetRecentClicks(ctx context.Context, shortCode, interval string) ([]entity.ClickByDate, error)
		GetClicksByBrowser(ctx context.Context, shortCode string) ([]entity.ClickByBrowser, error)
		GetClicksByDevice(ctx context.Context, shortCode string) ([]entity.ClickByDevice, error)
	}
)
