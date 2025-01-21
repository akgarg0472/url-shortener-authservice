# Step 1: Build the Go binary
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

ARG TARGETOS

ARG TARGETARCH

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-X 'github.com/akgarg0472/urlshortener-auth-service/build.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)' \
    -X 'github.com/akgarg0472/urlshortener-auth-service/build.GoVersion=$(go version | cut -d' ' -f3)' \
    -X 'github.com/akgarg0472/urlshortener-auth-service/build.OS=$(go env GOOS)' \
    -X 'github.com/akgarg0472/urlshortener-auth-service/build.Arch=$(go env GOARCH)' \
    -X 'github.com/akgarg0472/urlshortener-auth-service/build.AppVersion=$(cat VERSION)'" \
    -o authservice ./cmd/authservice/main.go

# Step 2: Create the final image to run the binary
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /app/authservice .

CMD ["./authservice"]
