CREATE TABLE IF NOT EXISTS urls
(
    id BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    short_code VARCHAR(50) UNIQUE NOT NULL,
    is_custom BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);