# Go User API

Production-oriented Gin + GORM user API with register, login, self-service profile actions, and admin-only user management.

## Endpoints

- `POST /api/register`
- `POST /api/login`
- `GET /api/me`
- `PUT /api/me`
- `GET /api/admin/users`
- `GET /api/admin/users/:id`
- `DELETE /api/admin/users/:id`

## Run locally

1. Set the environment variables in `.env`.
2. Start PostgreSQL with `docker compose up -d db`.
3. Run the API with `go run ./cmd`.
