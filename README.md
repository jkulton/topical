![topical-banner-lg](https://user-images.githubusercontent.com/6694167/112412673-4b1d0000-8cf5-11eb-9a38-bfa0d41a227e.png)

# Topical

[![Go Report Card](https://goreportcard.com/badge/github.com/jkulton/topical)](https://goreportcard.com/report/github.com/jkulton/topical)

Topical is a minimalist message board built with Golang.

Users can create topics and reply with mesages. Users need only to pick an avatar color and two initials to get a signature and start posting.

## Running Topical

1. Clone this repo
2. Navigate into the repo
3. Run `touch topical.db && go run ./scripts/db_reset` to create the database
4. Run `go run ./cmd/topical -session-key=secret` to start the app
5. Navigate to `localhost:8000` in your browser to use the app

### Running Tests

- Run tests by executing `go test ./internal/...`.
- You can also get coverage information with the `-cover` flag (`go test ./... -cover`)

### Options/Flags

Topical supports passing options via flag at start up or through an accompanying environment variable. Flags will take precedence over environment variables, if provided.

| Flag | ENV var | Default Fallback | Description |
|------|---------|---------|-------------|
| `p` | `PORT` | `8000`  | Port for Topical to bind to |
| `db-file` | `DB_FILE` | `'topical.db'` | Path to the sqlite3 database file |
| `session-key` | `SESSION_KEY` | `'not-set'` | Session key for cookie store |

### Database Management

A few DB management scripts have been provided and will accomplish the following tasks:

| Script | Use |
|--------|-----|
| `db_init` | Creates the intial tables in the DB as specified in `./schema.sql` |
| `db_seed` | Seeds an existing database with data from `./seeds.sql` |
| `db_drop` | Drops tables from a Topical database |
| `db_reset` | Drops tables, creates tables, reseeds database |

Scripts can be run directly with `go run`:

```sh
go run ./scripts/db_init -db-file='topical.db'
```
