# Chirpy Go Server

## Overview
This is a Go-based web server that serves static files, handles user authentication, manages chirps (short messages), and provides administrative metrics.

## Features
- Serves static files from `/app/`
- User authentication and JWT-based authorization
- Chirp management (CRUD operations)
- Webhook handling for Polka
- Health check endpoint
- Administrative metrics

## Requirements
- Go 1.18+
- PostgreSQL database
- `.env` file with necessary environment variables
- [Goose](https://github.com/pressly/goose) for database migrations
- [SQLC](https://sqlc.dev/) for query generation

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/your-repo.git
   cd your-repo
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Create a `.env` file and configure it:
   ```ini
   DB_URL=your_database_url
   JWT_SECRET=your_jwt_secret
   PLATFORM=your_platform_name
   POLKA_KEY=your_polka_api_key
   ```

## Database Setup
1. Run database migrations using Goose:
   ```sh
   goose -dir Chirpy/sql/schema postgres "$DB_URL" up
   ```
2. Generate queries using SQLC:
   ```sh
   sqlc generate
   ```

## Running the Server
```sh
  go run main.go
```

The server will start on port `8080` and serve static files from the current directory.

## API Endpoints

### Health Check
- `GET /api/healthz` - Check server readiness

### Authentication
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh JWT token
- `POST /api/revoke` - Revoke JWT token

### User Management
- `POST /api/users` - Create a user
- `PUT /api/users` - Update user details

### Chirp Management
- `GET /api/chirps` - Get all chirps
- `GET /api/chirps/{chirpID}` - Get a specific chirp
- `POST /api/chirps` - Create a chirp
- `DELETE /api/chirps/{chirpID}` - Delete a chirp

### Webhooks
- `POST /api/polka/webhooks` - Handle Polka webhook events

### Admin
- `GET /admin/metrics` - Get server metrics
- `POST /admin/reset` - Reset metrics

## License
This project is licensed under the MIT License.

