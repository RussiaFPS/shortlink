migrate_init:
	migrate create -ext sql -dir ./migrations -seq create_url_table

migrate_up:
	migrate -database "postgres://postgres:qwer1234@localhost:5432/shortlink?sslmode=disable" -path ./migrations up

migrate_down:
	migrate -database "postgres://postgres:qwer1234@localhost:5432/shortlink?sslmode=disable" -path ./migrations down

run:
	go run cmd/shortener/main.go -d "postgres://postgres:qwer1234@localhost:5432/shortlink"