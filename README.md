# URL Shortener Authentication Service

![Java Version](https://img.shields.io/badge/golang-1.21-blue)
![version](https://img.shields.io/badge/version-1.7.2-blue)

This project is a URL Shortener Authentication Service written in Go. It handles authentication, user management, token-based security (JWT), and integrates with other services like Kafka for email notifications and Consul for service discovery. This service is part of a larger URL shortener platform.

## Features

- **User Authentication**: JWT-based token authentication for secure API access.
- **Service Discovery**: Integration with Discovery Server for service registration and health checks.
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

- `ENABLE_DISCOVERY_CLIENT`: Enable Discovery client for service discovery (`true`/`false`). Default: `true`
- `DISCOVERY_SERVER_IP`: Discovery service URL for registering and heartbeat. Default: `http://localhost:8500`
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
- **Discovery Server** (optional, if you are using Consul/Eureka or any other for service discovery).

## Running the Project

### 1. Clone the Repository

```bash
git clone https://github.com/akgarg0472/url-shortener-authservice
cd url-shortener-authservice
```

### 2. Set Up Environment Variables

Create a `.env` file in the root directory of the project, and set the necessary environment variables as described above. Example `.env`:

```bash
LOGGING_CONSOLE_ENABLED=true
LOGGING_FILE_ENABLED=false
LOGGING_STREAM_ENABLED=false

# File Logging Settings
LOGGING_FILE_BASE_PATH=/tmp
LOG_LEVEL=INFO

# Stream Logging Settings
LOGGING_STREAM_HOST=localhost
LOGGING_STREAM_PORT=5000
LOGGING_STREAM_PROTOCOL=TCP

ENABLE_DISCOVERY_CLIENT=true
DISCOVERY_SERVER_IP=http://localhost:8500
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

## OAuth Provider Configuration

To enable OAuth authentication for your application, you'll need to configure the OAuth provider settings in your database. The following steps outline how to populate the OAuth provider table with the necessary details for GitHub and Google OAuth integration.

### 1. Add OAuth Providers to the Database:

Use the MySQL queries below to insert OAuth provider information into your `oauth_providers` table.

#### GitHub OAuth Configuration

To integrate GitHub OAuth, run the following SQL query:

```sql
INSERT INTO oauth_providers (provider, client_id, base_url, redirect_uri, access_type, scope)
VALUES
  ('github',
   'xxxxxxxx',
   'https://github.com/login/oauth/authorize',
   'http://localhost:3000/oauth/github/success',
   '',
   'user');
```

- **client_id**: Replace `'xxxxxxxx'` with your GitHub OAuth application's client ID.
- **redirect_uri**: Update with the appropriate redirect URI for your application (this should match the URL configured in your GitHub OAuth settings).
- **scope**: `'user'` specifies the permission scope for the user profile. Adjust as necessary for your app.

#### Google OAuth Configuration

To integrate Google OAuth, use this SQL query:

```sql
INSERT INTO oauth_providers (provider, client_id, base_url, redirect_uri, access_type, scope)
VALUES
  ('google',
   'xxxxxxxxxxxxx-yyyyyyyyyyyyyy.apps.googleusercontent.com',
   'https://accounts.google.com/o/oauth2/v2/auth',
   'http://localhost:3000/oauth/google/success',
   '',
   'openid email profile');
```

- **client_id**: Replace `'xxxxxxxxxxxxx-yyyyyyyyyyyyyy.apps.googleusercontent.com'` with your Google OAuth application's client ID.
- **redirect_uri**: Update with the appropriate redirect URI for your application (this should match the URL configured in your Google Cloud Console OAuth settings).
- **scope**: `'openid email profile'` grants access to the user’s basic profile, email, and OpenID information. Adjust the scope as per your application's requirements.

### 2. Restart the Authentication Service

After inserting the configuration into the table, **restart the authentication service** to apply the changes and enable OAuth functionality.

> ​
> **Note**:
>
> 1.  Ensure that you’ve properly registered your OAuth applications on GitHub and Google, and that your client IDs and secret keys are correctly configured.
> 2.  The provided redirect URIs should match the ones registered in your OAuth application settings on both GitHub and Google.
>     ​

For more details, refer to the official documentations:

- [GitHub OAuth Documentation](https://docs.github.com/en/developers/apps/building-oauth-apps)
- [Authorizing OAuth Apps](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps)
- [Google OAuth Documentation](https://developers.google.com/identity/protocols/oauth2)
