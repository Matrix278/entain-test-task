FROM golang:1.19.2-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

FROM migrate/migrate:v4.14.1

WORKDIR /app

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env ./.env
COPY --from=builder /app/main ./

# Set the environment variables
ENV POSTGRES_DATABASEURL=${POSTGRES_DATABASEURL}

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["sh", "-c", "source .env && migrate -path /app/migrations -database ${POSTGRES_DATABASEURL} -verbose up && ./main"]
