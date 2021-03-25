![Topical Banner](./assets/banner.png)

# Topical

A (very) minimalist message board

## Running Topical

There are a few usecases where you might want to run Topical, here they are:

- **Run Locally** (run Topical locally just to try it out)
- **Development** (you want to play with Topical's code and run it locally)
- **Deploy** (you want to deploy Topical for extended use)

### Run Locally

The easiest way to run Topical locally is using Docker. Please make sure you have both Docker and Docker Compose installed, then do the following from the command line:

1. Clone this repo
2. Navigate into the repo
3. Run `docker-compose up`
4. Navigate to `localhost:8000` in your browser. Done!

### Development

In order to edit Topical's code and then build the app please ensure you have Docker, Docker Compose, and Go installed. Then, do the following:

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

#### Running Tests

- Run tests by executing `go test ./...`.
- You can also get coverage information with the `-cover` flag (`go test ./... -cover`)

### Deploying to Heroku

(Heroku steps current as-of early 2021)

1. In the Heroku web interface create a new app (name it whatever you want, but remember it)
2. Navigate to the "Resources" tab for the app and the "Heroku Postgres" add-on (Free hobby plan is fine)
3. Navigate to the "Settings" tab for the app, click "Reveal Config Vars" and add a new var with the key `SESSION_KEY` and a "sufficiently random" value ([mux/sessions](https://pkg.go.dev/github.com/gorilla/sessions) wording)
4. Clone this repo locally
5. In the repo run `heroku git:remote -a YOUR_APP_NAME` where `YOUR_APP_NAME` is the name you chose in Step 1.
6. Commit and push the local code:
```
$ git add .
$ git commit -am "make it better"
$ git push heroku master
```
7. Wait for the Heroku build process to complete.
8. Navigate to the Heroku URL for your app

Note that the Docker Compose file in this directory should not be used for deployment on any platform as it defines the database and app session secret using dummy local values.
