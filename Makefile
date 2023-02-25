run:
	go run main.go

create-migration:
	@migrate create -ext sql -dir migrations -seq ${name}

migrate-up:
	@godotenv -f .env
	migrate -path migrations/ -database "${POSTGRES_DATABASEURL}" -verbose up

migrate-down:
	@godotenv -f .env
	@migrate -path migrations/ -database "${POSTGRES_DATABASEURL}" -verbose down
