# Hortons Diary

A privacy-conscious web application for recording cluster headache attacks and,
later, viewing personal summaries and charts.

Phase 2 provides a runnable Go application skeleton and PostgreSQL development
environment. Diary features, authentication and database migrations are not
implemented yet.

## Requirements

- Docker Engine with Docker Compose v2 for the complete local stack
- Go 1.27.0 or newer only when running the application directly on the host

## Run the complete stack

Copy `.env.example` to `.env` if you want to override the safe development
defaults, then run:

```sh
docker compose up --build
```

Open <http://localhost:8080>. Readiness is available at
<http://localhost:8080/health/ready>. Stop the stack with
`docker compose down`; the named PostgreSQL volume is retained. Use
`docker compose down --volumes` only when intentionally deleting local data.

## Run Go on the host

The application container is optional during development. Start only
PostgreSQL with a loopback-only port and then run Go directly:

```sh
docker compose -f compose.yaml -f compose.dev-db.yaml up -d db
set -a
. ./.env.example
set +a
go run ./cmd/server
```

The base Compose file does not expose PostgreSQL. The development override
publishes it only on `127.0.0.1`; do not use that override in production.
If port 5432 is already in use, choose another loopback port for both commands,
for example `POSTGRES_PORT=55432 docker compose -f compose.yaml -f
compose.dev-db.yaml up -d db` and set port `55432` in `DATABASE_URL`.

## Build and test

Unit and HTTP tests do not require Docker or a running database:

```sh
go test ./...
go build ./cmd/server
```

If the host does not have the pinned Go version, use the same toolchain as the
container build:

```sh
docker run --rm -v "$PWD:/src" -w /src golang:1.27.0-alpine go test ./...
```

Database integration tests will be added with the schema in phase 3. Those
tests will require PostgreSQL, but will not require the application container.

## Configuration

- `APP_ENV` (required): `development`, `test` or `production`
- `HTTP_ADDR`: listen address, defaults to `:8080`
- `DATABASE_URL` (required): PostgreSQL connection URL; never log or commit it
- `LOG_LEVEL`: `debug`, `info`, `warn` or `error`; defaults to `info`
- `SESSION_COOKIE_SECURE`: defaults to `true` in production and cannot be
  disabled there

The committed `.env.example` contains local-only values, not production
credentials.

## Documentation

- [Phase 1 architecture blueprint](docs/architecture.md)
- [ADR 0001: modular monolith with Go, HTMX and PostgreSQL](docs/adr/0001-modular-monolith-go-htmx-postgresql.md)
- [ADR 0002: initial Go application toolchain](docs/adr/0002-go-application-toolchain.md)
- [Product and deployment decisions](docs/decisions.md)

The phase 1 product and deployment decisions have been confirmed. The current
phase 2 skeleton follows the architecture blueprint without implementing phase
3 schema or later product functionality.
