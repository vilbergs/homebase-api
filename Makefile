migrate:
	docker run -v /Users/vilberg/go/src/github.com/vilbergs/homebase-api/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://homebase:homebase_pass@localhost:5432/homebase?sslmode=disable up
