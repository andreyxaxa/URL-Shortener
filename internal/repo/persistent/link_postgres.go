package persistent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/andreyxaxa/URL-Shortener/internal/entity"
	"github.com/andreyxaxa/URL-Shortener/pkg/postgres"
	"github.com/andreyxaxa/URL-Shortener/pkg/types/errs"
	"github.com/jackc/pgx/v5"
)

const (
	// Table
	urlsTable   = "urls"
	clicksTable = "clicks"

	// Column
	idColumn        = "id"
	urlColumn       = "url"
	shortCodeColumn = "short_code"
	isCustomColumn  = "is_custom"
	createdAtColumn = "created_at"

	urlIdColumn         = "url_id"
	ipAddrColumn        = "ip_address"
	userAgentColumn     = "user_agent"
	deviceColumn        = "device"
	browserFamilyColumn = "browser_family"
	clickedAtColumn     = "clicked_at"
)

type LinkRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *LinkRepo {
	return &LinkRepo{pg}
}

func (r *LinkRepo) GetNextSequenceValue(ctx context.Context) (int64, error) {
	sql, args, err := r.Builder.
		Select("nextval('urls_id_seq')").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("LinkRepo - GetNextSequenceValue - r.Builder.ToSql: %w", err)
	}

	var ID int64

	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&ID)
	if err != nil {
		return 0, fmt.Errorf("LinkRepo - GetNextSequenceValue - row.Scan: %w", err)
	}

	return ID, nil
}

func (r *LinkRepo) CreateWithShortCode(ctx context.Context, ID int64, originalURL, shortCode string, isCustom bool) error {
	if isCustom {
		sql, args, err := r.Builder.
			Insert(urlsTable).
			Columns(urlColumn, shortCodeColumn, isCustomColumn).
			Values(originalURL, shortCode, isCustom).
			ToSql()
		if err != nil {
			return fmt.Errorf("LinkRepo - CreateWithShortCode - r.Builder.ToSql: %w", err)
		}

		_, err = r.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("LinkRepo - CreateWithShortCode - r.Pool.Exec: %w", err)
		}

	} else {
		sql, args, err := r.Builder.
			Insert(urlsTable).
			Columns(idColumn, urlColumn, shortCodeColumn, isCustomColumn).
			Values(ID, originalURL, shortCode, isCustom).
			ToSql()
		if err != nil {
			return fmt.Errorf("LinkRepo - CreateWithShortCode - r.Builder.ToSql: %w", err)
		}

		_, err = r.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("LinkRepo - CreateWithShortCode - r.Pool.Exec: %w", err)
		}
	}

	return nil
}

func (r *LinkRepo) GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error) {
	sql, args, err := r.Builder.
		Select(urlColumn).
		From(urlsTable).
		Where(squirrel.Eq{shortCodeColumn: shortCode}).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("LinkRepo - GetOriginalURLByShortCode - r.Builder.ToSql: %w", err)
	}

	var url string

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("LinkRepo - GetOriginalURLByShortCode: %w", errs.ErrRecordNotFound)
		}
		return "", fmt.Errorf("LinkRepo - GetOriginalURLByShortCode - row.Scan: %w", err)
	}

	return url, nil
}

func (r *LinkRepo) GetIDByShortCode(ctx context.Context, shortCode string) (int64, error) {
	sql, args, err := r.Builder.
		Select(idColumn).
		From(urlsTable).
		Where(squirrel.Eq{shortCodeColumn: shortCode}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("LinkRepo - GetIDByShortCode - r.Builder.ToSql: %w", err)
	}

	var ID int64

	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("LinkRepo - GetIDByShortCode: %w", errs.ErrRecordNotFound)
		}
		return 0, fmt.Errorf("LinkRepo - GetIDByShortCode - row.Scan: %w", err)
	}

	return ID, nil
}

func (r *LinkRepo) CreateClick(ctx context.Context, urlID int64, IP, userAgent, device, browser string) error {
	sql, args, err := r.Builder.
		Insert(clicksTable).
		Columns(urlIdColumn, ipAddrColumn, userAgentColumn, browserFamilyColumn, deviceColumn, clickedAtColumn).
		Values(urlID, IP, userAgent, browser, device, time.Now()).
		ToSql()
	if err != nil {
		return fmt.Errorf("LinkRepo - CreateClick - r.Builder.ToSql: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("LinkRepo - CreateClick - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *LinkRepo) GetAnalytics(ctx context.Context, shortCode string) (entity.Analytics, error) {
	totalClicks, err := r.getTotalClicks(ctx, shortCode)
	if err != nil {
		return entity.Analytics{}, fmt.Errorf("LinkRepo - GetAnalytics - r.getTotalClicks: %w", err)
	}

	clicksByBrowser, err := r.GetClicksByBrowser(ctx, shortCode)
	if err != nil {
		return entity.Analytics{}, fmt.Errorf("LinkRepo - GetAnalytics - r.getClicksByBrowser: %w", err)
	}

	clicksByDevice, err := r.GetClicksByDevice(ctx, shortCode)
	if err != nil {
		return entity.Analytics{}, fmt.Errorf("LinkRepo - GetAnalytics - r.getClicksByDevice: %w", err)
	}

	// if we want full analytics - interval == day by default
	recentClicks, err := r.GetRecentClicks(ctx, shortCode, "day")
	if err != nil {
		return entity.Analytics{}, fmt.Errorf("LinkRepo - GetAnalytics - r.getRecentClicks: %w", err)
	}

	return entity.Analytics{
		TotalClicks:     totalClicks,
		ClicksByBrowser: clicksByBrowser,
		ClicksByDevice:  clicksByDevice,
		RecentClicks:    recentClicks,
	}, nil
}

func (r *LinkRepo) getTotalClicks(ctx context.Context, shortCode string) (int64, error) {
	sql := `
	SELECT COUNT(*) AS total_clicks
	FROM clicks c
	JOIN urls u ON u.id = c.url_id
	WHERE u.short_code = $1;
	`

	var total int64

	row := r.Pool.QueryRow(ctx, sql, shortCode)
	err := row.Scan(&total)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("LinkRepo - getTotalClicks: %w", errs.ErrRecordNotFound)
		}
		return 0, fmt.Errorf("LinkRepo - getTotalClicks - row.Scan: %w", err)
	}

	return total, nil
}

func (r *LinkRepo) GetClicksByBrowser(ctx context.Context, shortCode string) ([]entity.ClickByBrowser, error) {
	sql := `
	SELECT 
		c.browser_family, 
		COUNT (*) AS clicks
	FROM clicks c
	JOIN urls u ON u.id = c.url_id
	WHERE u.short_code = $1
	GROUP by c.browser_family
	ORDER BY clicks DESC;
	`

	rows, err := r.Pool.Query(ctx, sql, shortCode)
	if err != nil {
		return nil, fmt.Errorf("LinkRepo - getClicksByBrowser - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	clicks := make([]entity.ClickByBrowser, 0)

	for rows.Next() {
		var c entity.ClickByBrowser
		if err := rows.Scan(
			&c.Browser,
			&c.Clicks,
		); err != nil {
			return nil, fmt.Errorf("LinkRepo - getClicksByBrowser - rows.Next: %w", err)
		}
		clicks = append(clicks, c)
	}

	return clicks, nil
}

func (r *LinkRepo) GetClicksByDevice(ctx context.Context, shortCode string) ([]entity.ClickByDevice, error) {
	sql := `
	SELECT
		c.device,
		COUNT (*) AS clicks
	FROM clicks c
	JOIN urls u ON u.id = c.url_id
	WHERE u.short_code = $1
	GROUP BY c.device
	ORDER BY clicks DESC;
	`

	rows, err := r.Pool.Query(ctx, sql, shortCode)
	if err != nil {
		return nil, fmt.Errorf("LinkRepo - getClicksByDevice - r.Pool.Query: %w", err)
	}

	clicks := make([]entity.ClickByDevice, 0)

	for rows.Next() {
		var c entity.ClickByDevice
		if err := rows.Scan(
			&c.Device,
			&c.Clicks,
		); err != nil {
			return nil, fmt.Errorf("LinkRepo - getClicksByDevice - rows.Scan: %w", err)
		}
		clicks = append(clicks, c)
	}

	return clicks, nil
}

func (r *LinkRepo) GetRecentClicks(ctx context.Context, shortCode, interval string) ([]entity.ClickByDate, error) {
	sql := `
	SELECT
		date_trunc($2, c.clicked_at) AS click_date,
		COUNT (*) AS clicks
	FROM clicks c
	JOIN urls u ON u.id = c.url_id
	WHERE u.short_code = $1
	GROUP BY click_date
	ORDER BY click_date
	LIMIT 90;
	`

	rows, err := r.Pool.Query(ctx, sql, shortCode, interval)
	if err != nil {
		return nil, fmt.Errorf("LinkRepo - GetRecentClicks - r.Pool.Query: %w", err)
	}

	clicks := make([]entity.ClickByDate, 0)

	for rows.Next() {
		var c entity.ClickByDate
		if err := rows.Scan(
			&c.Date,
			&c.Clicks,
		); err != nil {
			return nil, fmt.Errorf("LinkRepo - GetRecentClicks - rows.Scan: %w", err)
		}
		clicks = append(clicks, c)
	}

	return clicks, nil
}

func (r *LinkRepo) ExistsByShortCode(ctx context.Context, shortCode string) error {
	sql, args, err := r.Builder.
		Select(idColumn).
		From(urlsTable).
		Where(squirrel.Eq{shortCodeColumn: shortCode}).
		ToSql()
	if err != nil {
		return fmt.Errorf("LinkRepo - ExistsByShortCode - r.Builder.ToSql: %w", err)
	}

	var d int

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&d); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrRecordNotFound
		}
		return fmt.Errorf("LinkRepo - ExistsByShortCode - r.Pool.QueryRow: %w", err)
	}

	return nil
}
