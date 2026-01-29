CREATE TABLE IF NOT EXISTS clicks 
(
    id BIGSERIAL PRIMARY KEY,
    url_id BIGINT REFERENCES urls(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT NOT NULL,
    device VARCHAR(30) NOT NULL,
    browser_family VARCHAR(50) NOT NULL,
    clicked_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_clicks_url_id ON clicks(url_id);
CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at);
CREATE INDEX IF NOT EXISTS idx_clicks_browser_family ON clicks(browser_family);
CREATE INDEX IF NOT EXISTS idx_clicks_device ON clicks(device);