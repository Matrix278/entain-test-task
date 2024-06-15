include .env

run:
	go run main.go

mod-vendor:
	go mod vendor

create-migration:
	@migrate create -ext sql -dir migrations -seq ${name}

migrate-up:
	@migrate -path ./migrations -database "${POSTGRES_DATABASEURL}" -verbose up

migrate-down:
	@migrate -path ./migrations -database "${POSTGRES_DATABASEURL}" -verbose down

migrate-force:
	@migrate -path ./migrations -database "${POSTGRES_DATABASEURL}" -verbose force ${version}
