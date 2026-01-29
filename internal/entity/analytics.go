package entity

import "time"

/*
{
	"analytics": {
    	"total_clicks": 1351,
    	"clicks_by_browser": [
      		{"browser": "Chrome", "clicks": 135},
      		{"browser": "Firefox", "clicks": 112},
      		{"browser": "Other", "clicks": 153}
    	],
    	"recent_clicks": [
      		{"date": "2026-01-27", "clicks": 12},
      		{"date": "2026-01-28", "clicks": 53}
    	]
  	}
}
*/

type Analytics struct {
	TotalClicks     int64            `json:"total_clicks"`
	ClicksByBrowser []ClickByBrowser `json:"clicks_by_browser"`
	ClicksByDevice  []ClickByDevice  `json:"clicks_by_device"`
	RecentClicks    []ClickByDate    `json:"recent_clicks"`
}

type ClickByDevice struct {
	Device string `json:"device"`
	Clicks int64  `json:"clicks"`
}

type ClickByBrowser struct {
	Browser string `json:"browser"`
	Clicks  int64  `json:"clicks"`
}

type ClickByDate struct {
	Date   time.Time `json:"date"`
	Clicks int64     `json:"clicks"`
}
