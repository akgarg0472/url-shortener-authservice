go build -ldflags "-X 'github.com/akgarg0472/urlshortener-auth-service/build.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)' \
-X 'github.com/akgarg0472/urlshortener-auth-service/build.GoVersion=$(go version | cut -d' ' -f3)' \
-X 'github.com/akgarg0472/urlshortener-auth-service/build.OS=$(go env GOOS)' \
-X 'github.com/akgarg0472/urlshortener-auth-service/build.Arch=$(go env GOARCH)' \
-X 'github.com/akgarg0472/urlshortener-auth-service/build.AppVersion=$(cat VERSION)'" \
-o authservice ./cmd/authservice/main.go

if [ $? -ne 0 ]; then
    echo "Build failed. Exiting."
    exit 1
fi

if [ -f ".env" ]; then
    set -a
    source .env
    set +a
else
    echo ".env file not found. Continuing with existing environment variables."
fi

clear
clear
./authservice
