CREATE TABLE urls (
                        id SERIAL PRIMARY KEY,
                        short_url VARCHAR(20) NOT NULL,
                        long_url  VARCHAR(255) NOT NULL,
                        CONSTRAINT unique_long_url UNIQUE (long_url),
                        CONSTRAINT unique_short_url UNIQUE (short_url)
);

CREATE INDEX idx_urls_shortURL ON urls(short_url);

CREATE INDEX idx_urls_longURL ON urls(long_url);