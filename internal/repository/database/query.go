package database

const addURL = `
		INSERT INTO urls (short_url, long_url, user_id)
			VALUES ($1, $2, $3)
		ON CONFLICT (long_url)
			DO UPDATE SET short_url = urls.short_url
		RETURNING short_url;
`

const findShortURL = `
		select short_url from urls
			where long_url = $1;
`

const findLongURL = `
		select long_url,is_deleted from urls
			where short_url = $1;
`

const selectFindAllByUser = `
	SELECT short_url, long_url FROM urls 
	    WHERE user_id = $1
`

const dellURL = `
	UPDATE urls SET is_deleted = true 
	    WHERE short_url = ANY($1) AND user_id = $2;
`

const createTable = `
	CREATE TABLE IF NOT EXISTS urls (
                        id SERIAL PRIMARY KEY,
                        short_url VARCHAR(255) NOT NULL,
                        long_url  VARCHAR(255) NOT NULL,
	    				user_id	VARCHAR(255),
	    				is_deleted BOOLEAN DEFAULT FALSE,
                        CONSTRAINT unique_long_url UNIQUE (long_url),
                        CONSTRAINT unique_short_url UNIQUE (short_url)
);

CREATE INDEX IF NOT EXISTS idx_urls_shortURL ON urls(short_url);

CREATE INDEX IF NOT EXISTS idx_urls_longURL ON urls(long_url);
`
