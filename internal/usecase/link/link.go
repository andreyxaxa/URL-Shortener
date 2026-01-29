package link

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andreyxaxa/URL-Shortener/internal/entity"
	"github.com/andreyxaxa/URL-Shortener/internal/repo"
	"github.com/andreyxaxa/URL-Shortener/pkg/encoder"
	"github.com/andreyxaxa/URL-Shortener/pkg/logger"
	"github.com/andreyxaxa/URL-Shortener/pkg/types/errs"
	"github.com/medama-io/go-useragent"
)

type LinkUseCase struct {
	repo  repo.LinkRepo
	cache repo.LinkCache

	logger logger.Interface
}

func New(r repo.LinkRepo, c repo.LinkCache, l logger.Interface) *LinkUseCase {
	return &LinkUseCase{
		repo:   r,
		cache:  c,
		logger: l,
	}
}

func (uc *LinkUseCase) CreateShortURL(ctx context.Context, originalURL string, customAlias string) (string, error) {
	var shortCode string
	var isCustom bool

	if customAlias != "" {
		err := uc.repo.ExistsByShortCode(ctx, customAlias)
		if err == nil {
			return "", fmt.Errorf("LinkUseCase - CreateShortURL: %w", errs.ErrAliasAlreadyTaken)
		}
		if !errors.Is(err, errs.ErrRecordNotFound) {
			return "", fmt.Errorf("LinkUseCase - CreateShortURL - uc.repo.ExistsByShortCode: %w", err)
		}

		shortCode = customAlias
		isCustom = true

		err = uc.repo.CreateWithShortCode(ctx, 0, originalURL, shortCode, isCustom)
		if err != nil {
			return "", fmt.Errorf("LinkUseCase - CreateShortURL - uc.repo.CreateWithShortCode: %w", err)
		}
	} else {
		nextID, err := uc.repo.GetNextSequenceValue(ctx)
		if err != nil {
			return "", fmt.Errorf("LinkUseCase - CreateShortURL - uc.repo.GetNextSequenceValue: %w", err)
		}

		shortCode = encoder.Encode(nextID)
		isCustom = false

		err = uc.repo.CreateWithShortCode(ctx, nextID, originalURL, shortCode, isCustom)
		if err != nil {
			return "", fmt.Errorf("LinkUseCase - CreateShortURL - uc.repo.CreateWithShortCode: %w", err)
		}
	}

	return shortCode, nil
}

func (uc *LinkUseCase) GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error) {
	cacheKey := fmt.Sprintf("url:%s", shortCode)

	// check cache
	originalURL, err := uc.cache.Get(ctx, cacheKey)
	if err == nil {
		// TODO: async with worker pool
		_, err = uc.cache.IncrementWithExpiry(ctx, fmt.Sprintf("hits:1h:%s", shortCode), 1*time.Hour)
		if err != nil {
			uc.logger.Warn("LinkUseCase - GetOriginalURLByShortCode - uc.cache.IncrementWithExpiry: %v", err)
		}

		_, err = uc.cache.IncrementWithExpiry(ctx, fmt.Sprintf("hits:24h:%s", shortCode), 24*time.Hour)
		if err != nil {
			uc.logger.Warn("LinkUseCase - GetOriginalURLByShortCode - uc.cache.IncrementWithExpiry: %v", err)
		}

		return originalURL, nil
	}

	if !errors.Is(err, errs.ErrRecordNotFound) {
		uc.logger.Warn("LinkUseCase - GetOriginalURLByShortCode - uc.cache.Get : %v", err)
	}

	// check repo
	originalURL, err = uc.repo.GetOriginalURLByShortCode(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("LinkUseCase - GetOriginalURLByShortCode - uc.repo.GetOriginalURLByShortCode: %w", err)
	}

	// cache set
	// TODO: async with worker pool
	ttl := uc.calculateTTL(ctx, shortCode)
	err = uc.cache.Set(ctx, cacheKey, originalURL, ttl)
	if err != nil {
		uc.logger.Warn("LinkUseCase - GetOriginalURLByShortCode - uc.cache.Set : %v", err)
	}

	return originalURL, nil
}

func (uc *LinkUseCase) calculateTTL(ctx context.Context, shortCode string) time.Duration {
	hits1h, err := uc.cache.GetInt(ctx, fmt.Sprintf("hits:1h:%s", shortCode))
	if err != nil {
		if !errors.Is(err, errs.ErrRecordNotFound) {
			uc.logger.Warn("LinkUseCase - calculateTTL - uc.cache.GetInt: %v", err)
		}
		hits1h = 0
	}

	hits24h, err := uc.cache.GetInt(ctx, fmt.Sprintf("hits:24h:%s", shortCode))
	if err != nil {
		if !errors.Is(err, errs.ErrRecordNotFound) {
			uc.logger.Warn("LinkUseCase - calculateTTL - uc.cache.GetInt: %v", err)
		}
		hits24h = 0
	}

	switch {
	case hits1h >= 100:
		return 3 * time.Hour
	case hits1h >= 20:
		return 1 * time.Hour
	case hits1h >= 5:
		return 30 * time.Minute
	case hits24h >= 50:
		return 15 * time.Minute
	case hits24h >= 10:
		return 10 * time.Minute
	default:
		return 5 * time.Minute
	}
}

// TODO: async with worker pool
func (uc *LinkUseCase) TrackClick(ctx context.Context, shortCode, IP, userAgent string) error {
	urlID, err := uc.repo.GetIDByShortCode(ctx, shortCode)
	if err != nil {
		return fmt.Errorf("LinkUseCase - TrackClick - uc.repo.GetIDByShortCode: %w", err)
	}

	// parse user-agent
	up := useragent.NewParser()
	agent := up.Parse(userAgent)
	device := agent.Device()
	browser := agent.Browser()

	err = uc.repo.CreateClick(ctx, urlID, IP, userAgent, device.String(), browser.String())
	if err != nil {
		return fmt.Errorf("LinkUseCase - TrackClick - uc.repo.CreateClick: %w", err)
	}

	return nil
}

func (uc *LinkUseCase) ExistsByShortCode(ctx context.Context, shortCode string) error {
	err := uc.repo.ExistsByShortCode(ctx, shortCode)
	if err != nil {
		return fmt.Errorf("LinkUseCase - ExistsByShortCode - uc.repo.ExistsByShortCode: %w", err)
	}

	return nil
}

func (uc *LinkUseCase) GetAnalytics(ctx context.Context, shortCode string) (entity.Analytics, error) {
	analytics, err := uc.repo.GetAnalytics(ctx, shortCode)
	if err != nil {
		return entity.Analytics{}, fmt.Errorf("LinkUseCase - GetAnalytics - uc.repo.GetAnalytics: %w", err)
	}

	return analytics, nil
}

func (uc *LinkUseCase) GetRecentClicks(ctx context.Context, shortCode, interval string) ([]entity.ClickByDate, error) {
	if interval != "day" && interval != "month" {
		return nil, fmt.Errorf("LinkUseCase - GetRecentClicks: %w", errs.ErrInvalidInterval)
	}

	analytics, err := uc.repo.GetRecentClicks(ctx, shortCode, interval)
	if err != nil {
		return nil, fmt.Errorf("LinkUseCase - GetRecentClicks - uc.repo.GetRecentClicks: %w", err)
	}

	return analytics, nil
}

func (uc *LinkUseCase) GetClicksByBrowser(ctx context.Context, shortCode string) ([]entity.ClickByBrowser, error) {
	analytics, err := uc.repo.GetClicksByBrowser(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("LinkUseCase - GetClicksByBrowser - uc.repo.GetClicksByBrowser: %w", err)
	}

	return analytics, nil
}

func (uc *LinkUseCase) GetClicksByDevice(ctx context.Context, shortCode string) ([]entity.ClickByDevice, error) {
	analytics, err := uc.repo.GetClicksByDevice(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("LinkUseCase - GetClicksByDevice - uc.repo.GetClicksByDevice: %w", err)
	}

	return analytics, nil
}
