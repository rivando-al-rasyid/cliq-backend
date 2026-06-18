# Cliq Backend

Cliq Backend is a REST API for a short link application. It handles user authentication, HttpOnly cookie sessions, profile management, password reset, short link creation, dashboard data, soft delete, and public slug redirection.

> This repository contains only the backend API. The React frontend should live in a separate repository and connect through `VITE_API_URL`.

---

## Features

- User registration and login
- JWT authentication with HttpOnly `access_token` cookie
- Bearer token support for API testing tools
- Current user endpoint
- Logout with token revocation and cookie clearing
- Password reset token flow
- Profile info and profile editing
- Avatar upload support
- Create short links with optional custom slug
- Reserved slug validation
- Dashboard links with pagination
- Soft delete links
- Public short link redirect by slug
- PostgreSQL migrations
- Redis connection support
- Swagger API documentation
- Docker Compose setup with PostgreSQL, Redis, pgAdmin, and migration service

---

## Tech Stack

- Go `1.26.3`
- Gin
- PostgreSQL `17`
- Redis `7`
- pgx
- JWT
- HttpOnly cookies
- golang-migrate `v4.18.3`
- Swagger
- Docker and Docker Compose

---

## Project Structure

```txt
cliq-backend/
├── cmd/
│   └── main.go
├── database/
│   └── migrations/
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internals/
│   ├── cache/
│   ├── config/
│   ├── controller/
│   ├── dto/
│   ├── middleware/
│   ├── model/
│   ├── pkg/
│   ├── repository/
│   ├── router/
│   └── service/
├── public/
│   └── img/
├── Dockerfile
├── docker-compose.yml
├── env.example
├── go.mod
├── go.sum
└── Makefile
```

---

## Requirements

For Docker development:

- Docker
- Docker Compose
- Make

For local development without Docker:

- Go `1.26.3`
- PostgreSQL
- Redis
- golang-migrate
- Make

---

## Environment Variables

Copy the example file:

```bash
cp env.example .env
```

Example `.env` for Docker development:

```env
# App
APP_NAME=Cliq
APP_ENV=development
APP_PORT=8080

# Database
DB_USER=cliq
DB_PASS=cliq_password
DB_NAME=cliq_db
DB_HOST=postgres
DB_PORT=5432

# JWT
JWT_ISSUER=cliq
JWT_SECRET=change_this_to_a_long_random_secret
JWT_EXPIRED=15m

# Redis
RDB_ADDR=redis:6379
RDB_USER=
RDB_PASS=

# Cookie and CORS
COOKIE_SECURE=false
COOKIE_SAMESITE=lax
ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173

# Short link response base URL
SHORT_LINK_BASE_URL=http://localhost:8080
```

For local development without Docker services, change these values:

```env
DB_HOST=localhost
RDB_ADDR=localhost:6379
```

Important notes:

- Keep `.env` out of Git.
- Use a long random value for `JWT_SECRET`.
- The current server starts on `0.0.0.0:8080` in `cmd/main.go`. `APP_PORT` exists in the env file, but the current code does not read it yet.
- For cross-origin frontend requests, keep `ALLOWED_ORIGINS` in sync with the frontend URL.

---

## Run with Docker

Start the backend stack:

```bash
docker compose up -d --build
```

Run migrations:

```bash
make migrate-up
```

Backend API:

```txt
http://localhost:8080
```

Swagger documentation:

```txt
http://localhost:8080/swagger/index.html
```

pgAdmin:

```txt
http://localhost
```

Stop containers:

```bash
docker compose down
```

Stop containers and remove volumes:

```bash
docker compose down -v
```

---

## Run Locally

Install dependencies:

```bash
go mod download
```

Make sure PostgreSQL and Redis are running, then run migrations:

```bash
make migrate-up
```

Start the API:

```bash
go run ./cmd
```

The API runs at:

```txt
http://localhost:8080
```

---

## Migration Commands

```bash
# create a new migration
make migrate-create NAME=example

# run all pending migrations
make migrate-up

# rollback migrations
make migrate-down

# force migration version
make migrate-force VERSION=1

# show current migration version
make migrate-status
```

The migration service uses the `migrate/migrate:v4.18.3` Docker image from `docker-compose.yml`.

---

## Authentication

Cliq uses a short-lived JWT access token.

For browser clients:

1. The client sends credentials to `POST /auth/login`.
2. The API validates the user.
3. The API sets an HttpOnly cookie named `access_token`.
4. The browser sends the cookie automatically when the frontend uses `credentials: "include"`.
5. Protected middleware reads the cookie and validates the token.

For API testing tools, protected routes also accept:

```txt
Authorization: Bearer <token>
```

Token lifetime:

```txt
Access token: 15 minutes
Password reset JWT: 10 minutes
```

---

## API Routes

### Auth

| Method | Endpoint | Description | Auth |
| --- | --- | --- | --- |
| POST | `/auth/register` | Register a new user | Public |
| POST | `/auth/login` | Login and set HttpOnly cookie | Public |
| GET | `/auth/me` | Get current authenticated user | Cookie or Bearer |
| POST | `/auth/logout` | Revoke current token and clear cookie | Cookie or Bearer |
| POST | `/auth/reset` | Request password reset token | Public |
| POST | `/auth/reset/confirm` | Exchange reset token for reset JWT | Public |
| POST | `/auth/change-password` | Set new password using reset JWT | Bearer reset JWT |

### Links

| Method | Endpoint | Description | Auth |
| --- | --- | --- | --- |
| POST | `/link/create` | Create a short link | Cookie or Bearer |
| GET | `/link/dashboard?page=1&limit=10` | Get paginated dashboard links | Cookie or Bearer |
| DELETE | `/link/:id` | Soft delete a link | Cookie or Bearer |
| GET | `/:slug` | Redirect to the original URL | Public |

### Profile

| Method | Endpoint | Description | Auth |
| --- | --- | --- | --- |
| GET | `/profile/info` | Get compact user info | Cookie or Bearer |
| GET | `/profile` | Get full profile | Cookie or Bearer |
| PATCH | `/profile/edit` | Edit profile and upload photo | Cookie or Bearer |
| PATCH | `/profile/change/password` | Change password using old password | Cookie or Bearer |

### Static Files and Docs

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/img/*` | Serve uploaded images from `public/img` |
| GET | `/swagger/index.html` | Swagger API documentation |

---

## Request Examples

### Register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Login and Save Cookie

```bash
curl -i -c cookies.txt -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Get Current User with Cookie

```bash
curl -b cookies.txt http://localhost:8080/auth/me
```

### Create Short Link with Cookie

```bash
curl -b cookies.txt -X POST http://localhost:8080/link/create \
  -H "Content-Type: application/json" \
  -d '{
    "origin_link": "https://example.com/very/long/url",
    "slug": "my-link"
  }'
```

### Create Short Link Without Custom Slug

```bash
curl -b cookies.txt -X POST http://localhost:8080/link/create \
  -H "Content-Type: application/json" \
  -d '{
    "origin_link": "https://example.com/very/long/url"
  }'
```

### Dashboard Links

```bash
curl -b cookies.txt "http://localhost:8080/link/dashboard?page=1&limit=10"
```

### Redirect

```bash
curl -I http://localhost:8080/my-link
```

---

## API Response Format

Successful responses usually follow this shape:

```json
{
  "message": "Short link created successfully",
  "data": {},
  "isSuccess": true
}
```

Error responses usually follow this shape:

```json
{
  "message": "Invalid request payload",
  "isSuccess": false,
  "error": "error detail"
}
```

---

## Database Summary

The migrations create these main resources:

- `users`
- `profiles`
- `token_type` enum
- `tokens`
- `links`
- related indexes

The `tokens` table stores access tokens and reset tokens with expiration and revocation status. The `links` table stores original URLs, unique slugs, click counts, owner ID, and soft-delete fields.

Slug rules:

- Must be 3-50 characters
- Can contain letters, numbers, and hyphens
- Must be unique
- Should not conflict with reserved application paths such as `auth`, `dashboard`, `link`, `profile`, `swagger`, or `img`

---

## Frontend Integration

The frontend should call this API with credentials enabled.

Example frontend env:

```env
VITE_API_URL=http://localhost:8080
VITE_SHORT_URL_BASE=http://localhost:8080
```

Example fetch behavior:

```js
fetch(`${import.meta.env.VITE_API_URL}/auth/me`, {
  credentials: "include",
});
```

CORS must allow the frontend origin:

```env
ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
```

---

## Development Notes

- Do not store the access token in frontend `localStorage`.
- Use the HttpOnly cookie flow for browser authentication.
- Use Bearer tokens only for manual API testing or reset-password JWT flow.
- Use `COOKIE_SECURE=true` in production with HTTPS.
- Use `COOKIE_SAMESITE=none` only when frontend and backend are on different HTTPS sites and cross-site cookies are required.
- Keep Swagger comments synchronized with the actual router paths.

---

## License

This project is licensed under the MIT License.
