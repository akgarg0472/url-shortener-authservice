# Step 1: Build the Go binary
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o authservice ./cmd/authservice/main.go

# Step 2: Create the final image to run the binary
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /app/authservice .

CMD ["./authservice"]
