include .env

run:
	go run main.go

create-migration:
	@migrate create -ext sql -dir migrations -seq ${name}

migrate-up:
	@migrate -path ./migrations -database "${POSTGRES_DATABASEURL}" -verbose up

migrate-down:
	@migrate -path ./migrations -database "${POSTGRES_DATABASEURL}" -verbose down

migrate-force:
	@migrate -path ./migrations -database "${POSTGRES_DATABASEURL}" -verbose force ${version}

docker-build:
	docker build -t entain .

docker-run:
	docker run -p ${SERVER_PORT}:8080 --env-file=.env entain

docker-compose-build:
	docker-compose build

docker-compose-up:
	docker-compose up
