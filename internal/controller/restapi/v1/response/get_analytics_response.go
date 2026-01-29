package response

import "github.com/andreyxaxa/URL-Shortener/internal/entity"

// Full Analytics

type GetAnalyticsResponse struct {
	Analytics Analytics `json:"analytics"`
}

type Analytics struct {
	TotalClicks     int64                   `json:"total_clicks"`
	ClicksByBrowser []entity.ClickByBrowser `json:"clicks_by_browser"`
	ClicksByDevice  []entity.ClickByDevice  `json:"clicks_by_device"`
	RecentClicks    []ClickByDate           `json:"recent_clicks"`
}

type ClickByDate struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

// Analytics by day / month

type GetAnalyticsByDateResponse struct {
	Analytics AnalyticsByDate `json:"analytics"`
}

type AnalyticsByDate struct {
	RecentClicks []ClickByDate `json:"recent_clicks"`
}

// Analytics by browser

type GetAnalyticsByBrowserResponse struct {
	Analytics AnalyticsByBrowser `json:"analytics"`
}

type AnalyticsByBrowser struct {
	ClicksByBrowser []entity.ClickByBrowser `json:"clicks_by_browser"`
}

// Analytics by device

type GetAnalyticsByDeviceResponse struct {
	Analytics AnalyticsByDevice `json:"analytics"`
}

type AnalyticsByDevice struct {
	ClicksByDevice []entity.ClickByDevice `json:"clicks_by_device"`
}
