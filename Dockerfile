FROM golang:1.22.4-alpine3.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Install any needed packages specified in go.mod and build the go app
RUN apk update && \
    apk add --no-cache git bash build-base && \
    go mod download && \
    go build -o main .

# This container exposes port 8080 to the outside world
EXPOSE ${SERVER_PORT}
