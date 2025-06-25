# Hirelin Auth Service

A modular authentication service for the Hirelin platform, built in Go. Supports OAuth (Google, GitHub), session management, and secure password handling.

## Features

- OAuth 2.0 authentication (Google, GitHub)
- Secure session and refresh token management (JWT)
- Password hashing with bcrypt
- Email verification support
- PostgreSQL database integration (via [sqlc](https://sqlc.dev/))
- CORS and logging middleware

## Project Structure

- `cmd/` - Main entrypoint, OAuth logic, adapters
- `internal/` - Core logic (routes, server, logger, utils, CORS)
- `db/` - SQL queries for sqlc
- `.env.example` - Example environment variables

## Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL database
- [sqlc](https://sqlc.dev/) for generating Go DB code

### Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/your-org/hirelin-auth.git
   cd hirelin-auth
   ```

2. **Configure environment variables:**

   - Copy `.env.example` to `.env` and fill in the required values.

3. **Generate database code:**

   ```bash
   sqlc generate
   ```

4. **Run the server:**
   ```bash
   go run ./cmd/main.go
   ```

## Usage

- **Ping endpoint:**  
  `GET /api/ping`  
  Returns `{ "message": "pong" }` for health checks.

- **OAuth sign-in:**  
  `GET /api/auth/oauth/signin?provider=google&redirect=/dashboard`  
  Redirects to the provider's OAuth consent screen.

- **OAuth callback:**  
  Handled automatically at `/api/auth/callback/{provider}`.

- **Session cookies:**  
  On successful login, sets `session_id` and `refresh_jwt` cookies.

## Environment Variables

See `.env.example` for all required variables, including:

- `DATABASE_URL`
- `HOST`, `PORT`
- `CLIENT_URL`, `SERVER_URL`
- `AUTH_SECRET`, `JWT_SECRET`
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`

## Database

- Uses PostgreSQL.
- SQL queries are defined in `db/query.sql` and used via sqlc.

## Extending

- Add new OAuth providers by updating `cmd/oauth/constants.go` and `ProviderEndpoints`.
- Add new routes in `internal/routes`.

## License

MIT License.
