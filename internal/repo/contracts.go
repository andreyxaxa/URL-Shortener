package repo

import (
	"context"
	"time"

	"github.com/andreyxaxa/URL-Shortener/internal/entity"
)

type (
	LinkRepo interface {
		GetNextSequenceValue(ctx context.Context) (int64, error)
		CreateWithShortCode(ctx context.Context, ID int64, originalURL, shortCode string, isCustom bool) error
		GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error)
		GetIDByShortCode(ctx context.Context, shortCode string) (int64, error)
		CreateClick(ctx context.Context, urlID int64, IP, userAgent, device, browser string) error
		GetAnalytics(ctx context.Context, shortCode string) (entity.Analytics, error)
		GetRecentClicks(ctx context.Context, shortCode, interval string) ([]entity.ClickByDate, error)
		GetClicksByBrowser(ctx context.Context, shortCode string) ([]entity.ClickByBrowser, error)
		GetClicksByDevice(ctx context.Context, shortCode string) ([]entity.ClickByDevice, error)
		// ExistsByShortCode returns error if record not exists, nil if record exists
		ExistsByShortCode(ctx context.Context, shortCode string) error
	}

	LinkCache interface {
		Get(ctx context.Context, key string) (string, error)
		GetInt(ctx context.Context, key string) (int64, error)
		Set(ctx context.Context, key string, value string, ttl time.Duration) error
		Delete(ctx context.Context, key string) error
		Increment(ctx context.Context, key string) (int64, error)
		IncrementWithExpiry(ctx context.Context, key string, ttl time.Duration) (int64, error)
	}
)
