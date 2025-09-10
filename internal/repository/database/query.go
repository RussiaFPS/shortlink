package database

const addURL = `
		INSERT INTO urls (short_url, long_url)
			VALUES ($1, $2)
		ON CONFLICT (long_url)
			DO NOTHING
		RETURNING short_url;
`

const findShortURL = `
		select short_url from urls
			where long_url = $1;
`

const findLongURL = `
		select long_url from urls
			where short_url = $1;
`

const createTable = `
	CREATE TABLE IF NOT EXISTS urls (
                        id SERIAL PRIMARY KEY,
                        short_url VARCHAR(255) NOT NULL,
                        long_url  VARCHAR(255) NOT NULL,
                        CONSTRAINT unique_long_url UNIQUE (long_url),
                        CONSTRAINT unique_short_url UNIQUE (short_url)
);
`
