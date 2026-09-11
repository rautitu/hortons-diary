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

With the example environment, open <http://localhost:51201>. Readiness is
available at <http://localhost:51201/health/ready>. Stop the stack with
`docker compose down`; the named PostgreSQL volume is retained. Use
`docker compose down --volumes` only when intentionally deleting local data.

`APP_PORT` controls the application's loopback-only host port. PostgreSQL stays
inside the Compose network unless the development override described below is
used.

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

The base Compose file does not expose PostgreSQL. The development override uses
`POSTGRES_PORT` to publish it only on `127.0.0.1`; do not use that override in
production. `.env.example` uses port 51202 consistently for the override and
the host application's `DATABASE_URL`.

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
- `LOG_LEVEL`: `debug`, `info`, `warn` or `error`; defaults to `info`
- `SESSION_COOKIE_SECURE`: defaults to `true` in production and cannot be
  disabled there
- `APP_PORT`: loopback-only host port published by the app container
- `POSTGRES_PORT`: loopback-only host port published by the dev DB override
- `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`: development database
  name and credentials used by Compose
- `HTTP_ADDR`: listen address used when Go runs directly on the host
- `DATABASE_URL`: PostgreSQL URL used when Go runs directly on the host; never
  log or commit production credentials

Compose deliberately gives the app container its internal `:8080` listener and
`db:5432` database address. `APP_PORT` and `POSTGRES_PORT` map those internal
ports to the host, while `HTTP_ADDR` and `DATABASE_URL` are for host Go runs.

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
