# go-gator

## Setup

Create `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://example"
}
```

## Build & Run

You'll need Go >= 1.24.0 and Postgres to run this project.

For local DX a `docker-compose.yml` with Postgres is provided.

Build with `go build -o gator && ./gator` or run via `go run .` during development.

For a list of available commands see the `handler_*` files or `main.go` where commands are registered.
