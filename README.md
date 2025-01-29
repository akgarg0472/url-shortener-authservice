# URL Shortener Authentication Service

![Java Version](https://img.shields.io/badge/golang-1.21-blue)
![version](https://img.shields.io/badge/version-1.5.0-blue)

This project is a URL Shortener Authentication Service written in Go. It handles authentication, user management, token-based security (JWT), and integrates with other services like Kafka for email notifications and Eureka for service discovery. This service is part of a larger URL shortener platform.

## Features

- **User Authentication**: JWT-based token authentication for secure API access.
- **Service Discovery**: Integration with Eureka for service registration and health checks.
- **Database Integration**: MySQL database for storing user data and other relevant information.
- **Kafka Integration**: Publish email notifications to Kafka topics.
- **OAuth Integration**: Supports OAuth with Google and GitHub for user authentication.

## Environment Variables

The project relies on a set of environment variables for configuration. Below is a list of all available environment variables with their descriptions.

### Logging Configuration

- `LOGGER_LEVEL`: Set the logging level (`DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`). Default: `INFO`
- `LOGGER_TYPE`: The type of logger to use (`console`, `file`). Default: `console`
- `LOGGER_ENABLED`: Enable or disable logging. Default: `true`
- `LOGGER_LOG_FILE_PATH`: Path to the log file if `LOGGER_TYPE` is set to `file`. Default: `/tmp/logs.log`

### Service Discovery

- `ENABLE_DISCOVERY_CLIENT`: Enable Eureka client for service discovery (`true`/`false`). Default: `true`
- `DISCOVERY_CLIENT_IP`: Eureka service URL for registering and heartbeat. Default: `http://localhost:8761/eureka/v2`
- `DISCOVERY_CLIENT_HEARTBEAT_FREQEUENCY_DURATION`: Heartbeat frequency in seconds. Default: `30`
-

### MySQL Database Configuration

- `MYSQL_DB_USERNAME`: MySQL database username. Default: `root`
- `MYSQL_DB_PASSWORD`: MySQL database password. Default: `root`
- `MYSQL_DB_HOST`: MySQL host. Default: `127.0.0.1`
- `MYSQL_DB_PORT`: MySQL port. Default: `3306`
- `MYSQL_DB_NAME`: MySQL database name. Default: `urlshortener`
- `MYSQL_USERS_TABLE_NAME`: The name of the users table. Default: `users`

### MySQL Connection Pool

- `MYSQL_CONNECTION_POOL_MAX_IDLE_CONNECTION`: Maximum number of idle connections in the MySQL connection pool. Default: `5`
- `MYSQL_CONNECTION_POOL_MAX_OPEN_CONNECTION`: Maximum number of open connections in the MySQL connection pool. Default: `10`

### JWT Authentication Configuration

- `JWT_SECRET_KEY`: Secret key used to sign JWT tokens.
- `JWT_TOKEN_ISSUER`: The issuer of the JWT token. Default: `urlshortener-auth-service`
- `JWT_TOKEN_EXPIRY`: Expiry time of the JWT token in seconds. Default: `60000` (60 seconds)

### Kafka Integration

- `KAFKA_CONNECTION_URL`: Kafka connection URL. Default: `localhost:9092`
- `KAFKA_TOPIC_EMAIL_NOTIFICATION`: Kafka topic for email notifications. Default: `urlshortener.notifications.email`
- `KAFKA_TOPIC_USER_REGISTERED`: Kafka topic for user registration successful. Default: `user.registration.completed`

### Forgot Password Configuration

- `FORGOT_PASS_SECRET_KEY`: Secret key for verifying forgot password tokens.
- `FORGOT_PASS_EXPIRY`: Expiry time of the forgot password token in seconds. Default: `600` (10 minutes)

### Frontend & Backend Configuration

- `BACKEND_BASE_DOMAIN`: Base URL for the backend API. Default: `http://localhost:8765/`
- `BACKEND_RESET_PASSWORD_URL`: API endpoint for resetting the password. Default: `api/v1/auth/verify-reset-password`
- `FRONTEND_BASE_DOMAIN`: Base URL for the front-end application. Default: `http://127.0.0.1:3000/`
- `FRONTEND_RESET_PASSWORD_PAGE_URL`: URL path for the reset password page. Default: `reset-password`
- `FRONTEND_DASHBOARD_PAGE_URL`: URL path for the dashboard page. Default: `dashboard`

### URL Shortener Configuration

- `URL_SHORTENER_LOGO_URL`: URL for the logo of the URL shortener platform.

### OAuth Configuration

- `OAUTH_GOOGLE_CLIENT_ID`: Google OAuth client ID for user authentication.
- `OAUTH_GOOGLE_CLIENT_SECRET`: Google OAuth client secret for user authentication.
- `OAUTH_GOOGLE_CLIENT_REDIRECT_URI`: Redirect URI for Google OAuth success (Front-end).
- `OAUTH_GITHUB_CLIENT_ID`: GitHub OAuth client ID for user authentication.
- `OAUTH_GITHUB_CLIENT_SECRET`: GitHub OAuth client secret for user authentication.
- `OAUTH_GITHUB_CLIENT_REDIRECT_URI`: Redirect URI for GitHub OAuth success (Front-end).

## Prerequisites

Make sure you have the following installed on your system:

- **Go (1.21 or higher)**: This project is built with Go.
- **Docker** (optional, if you want to use Docker to run the service).
- **MySQL**: For local development, you need a MySQL database running.
- **Kafka** (optional, if you want to use Kafka for notifications).
- **Eureka Server** (optional, if you are using Eureka for service discovery).

## Running the Project

### 1. Clone the Repository

```bash
git clone https://github.com/akgarg0472/url-shortener-authservice
cd url-shortener-authservice
```

### 2. Set Up Environment Variables

Create a `.env` file in the root directory of the project, and set the necessary environment variables as described above. Example `.env`:

```bash
LOGGER_LEVEL=INFO
LOGGER_TYPE=console
LOGGER_ENABLED=true
LOGGER_LOG_FILE_PATH=/tmp/logs.log

ENABLE_DISCOVERY_CLIENT=true
DISCOVERY_CLIENT_IP=http://localhost:8761/eureka/v2
DISCOVERY_CLIENT_HEARTBEAT_FREQEUENCY_DURATION=30

MYSQL_DB_USERNAME=root
MYSQL_DB_PASSWORD=root
MYSQL_DB_HOST=127.0.0.1
MYSQL_DB_PORT=3306
MYSQL_DB_NAME=urlshortener
MYSQL_USERS_TABLE_NAME=users

MYSQL_CONNECTION_POOL_MAX_IDLE_CONNECTION=5
MYSQL_CONNECTION_POOL_MAX_OPEN_CONNECTION=10

JWT_SECRET_KEY=KFwdkp3zZnGx89LSukJpFR2rRjk7zm
JWT_TOKEN_ISSUER=urlshortener-auth-service
JWT_TOKEN_EXPIRY=60000

KAFKA_CONNECTION_URL=localhost:9092
KAFKA_TOPIC_EMAIL_NOTIFICATION=urlshortener.notifications.email
KAFKA_TOPIC_USER_REGISTERED=user.registration.completed

FORGOT_PASS_SECRET_KEY=z8MTCxxKUvHCgQ9rgCP9Si50haCa6y
FORGOT_PASS_EXPIRY=600

BACKEND_BASE_DOMAIN=http://localhost:8765/
BACKEND_RESET_PASSWORD_URL=api/v1/auth/verify-reset-password

FRONTEND_BASE_DOMAIN=http://127.0.0.1:3000/
FRONTEND_RESET_PASSWORD_PAGE_URL=reset-password
FRONTEND_DASHBOARD_PAGE_URL=dashboard

URL_SHORTENER_LOGO_URL=

OAUTH_GOOGLE_CLIENT_ID=your-google-client-id
OAUTH_GOOGLE_CLIENT_SECRET=your-google-client-secret
OAUTH_GOOGLE_CLIENT_REDIRECT_URI=http://localhost:3000/oauth/google/success

OAUTH_GITHUB_CLIENT_ID=your-github-client-id
OAUTH_GITHUB_CLIENT_SECRET=your-github-client-secret
OAUTH_GITHUB_CLIENT_REDIRECT_URI=http://localhost:3000/oauth/github/success
```

### 3. Build the Project

```bash
go build -o authservice ./cmd/authservice/main.go
```

### 4. Run the Project

```bash
./authservice
```

## Docker Setup

The application is Dockerized for simplified deployment. The `Dockerfile` is already configured to build and run the
Spring Boot application.

The `Dockerfile` defines the build and runtime configuration for the container.

### Building the Docker Image

To build the Docker image, run the following command:

```bash
docker build -t akgarg0472/urlshortener-auth-service:1.0.0 .
```

### Run the Docker Container

You can run the application with custom environment variables using the docker run command. For example:

```bash
docker run --network=host --env-file .env akgarg0472/urlshortener-auth-service:1.0.0
```
