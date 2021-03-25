![topical-banner-lg](https://user-images.githubusercontent.com/6694167/112412673-4b1d0000-8cf5-11eb-9a38-bfa0d41a227e.png)

# Topical

[Live Demo](https://topical-go.herokuapp.com/)

Topical is a (very) minimalist message board built with Golang.

Users can create topics and reply with mesages. Users need only to pick an avatar color and two initials to get a signature and start posting.

Topical even supports **dark mode**. 😍

Topical uses the preferred color scheme from your OS settings to decide which theme to display (update your OS setting to see!)

## Running Topical

The simplest way to run Topical is using Docker. Please ensure you have both Docker and Docker Compose installed, then do the following:

1. Clone this repo
2. Navigate into the repo
3. Run `docker-compose up`
4. Navigate to `localhost:8000` in your browser. Done!

---

## Development

Please ensure you have both Docker and Docker Compose installed, then do the following:

1. Clone this repo
2. Navigate into the repo
3. Run `docker-compose up db`
4. In another terminal window, in the repo directory also run:

```
go run ./cmd/topical \
  -db-connection-uri='postgresql://topical:topical@localhost?sslmode=disable' \
  -session-key=somethingSecret
```

5. Navigate to `localhost:8000` in your browser. Done!

You can change Topical's code locally and re-run the `go run` command to restart the server.

### Running Tests

- Run tests by executing `go test ./...`.
- You can also get coverage information with the `-cover` flag (`go test ./... -cover`)

### Options/Flags

Topical supports passing options via flag at start up or through an accompanying environment variable. Flags will take precedence over environment variables, if provided.

| Flag | ENV var | Default Fallback | Description |
|------|---------|---------|-------------|
| `p` | `PORT` | `8000`  | Port for Topical to bind to |
| `database-url` | `DATABASE_URL` | `'not-set'` | URI-formatted Postgres connection information (e.g. `postgresql://localhost:5433...`) |
| `session-key` | `SESSION_KEY` | `'not-set'` | Session key for cookie store |

### Database Management

A few DB management scripts have been provided and will accomplish the following tasks:

| Script | Use |
|--------|-----|
| `db_init` | Creates the intial tables in the DB as specified in `./schema.sql` |
| `db_seed` | Seeds an existing database with data from `./seeds.sql` |
| `db_drop` | Drops tables from a Topical database |
| `db_reset` | Drops tables, creates tables, reseeds database |

Scripts can be run directly with `go run` but require passing a `-database-url` (if the environment variable is not set):

```sh
go run ./scripts/db_init \
  -database-url='postgresql://topical:topical@localhost?sslmode=disable'
```
