package repository

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
