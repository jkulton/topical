![Topical Banner](./assets/banner.png)

# Topical

[Live Demo](https://topical-go.herokuapp.com/)

Topical is a (very) minimalist message board built with Golang and a few external dependencies.

Users can post topics and reply with mesages. Users "sign up" by storing a cookie associated with their favorite color and two initials in their browser. The user can then post on the message board and their messages will be tied to this signature. The messages are persisted to a Postgres database.

## Running Locally

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
4. In another terminal window, also in the repo's directory run:

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

Topical supports passing options via flag at start up or through an accompanying environment variable. Flag will take precedence over the environment variables, if provided.

| Flag | ENV var | Default Fallback | Description |
|------|---------|---------|-------------|
| `-p` | `PORT` | `8000`  | Port for Topical to bind to |
| `-database-url` | `DATABASE_URL` | `'not-set'` | URI-formatted Postgres connection information (e.g. `postgresql://localhost:5433...`) |
| `-session-key` | `SESSION_KEY` | `'not-set'` | session key for cookie store |

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
