![Topical Banner](./assets/banner.png)

# Topical

A (very) minimalist message board

## Run Locally

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

---

## Deploying to Heroku

(Heroku steps current as-of early 2021)

1. In the Heroku web interface create a new app
2. Navigate to the "Resources" tab for the app and add the "Heroku Postgres" add-on (Free Hobby plan works fine)
3. Navigate to the "Settings" tab for the app, click "Reveal Config Vars" and add a new var with the key `SESSION_KEY` and a "sufficiently random" value (see [mux/sessions](https://pkg.go.dev/github.com/gorilla/sessions) for what this means)
4. Clone this repo locally
5. In the repo run `heroku git:remote -a YOUR_APP_NAME` where `YOUR_APP_NAME` is the name you chose in Step 1.
6. Commit and push the local code:
```sh
$ git add .
$ git commit -am "make it better"
$ git push heroku master
```
7. Wait for the Heroku build process to complete.
8. Navigate to the Heroku URL for your app. Done!

Note that the Docker Compose file in this directory should not be used for deployment on any platform as it defines the database and app session secret using dummy local values.
