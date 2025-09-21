# Blog aggreGATOR project from Boot.Dev

Simple blog aggre**GATOR** project from Boot.Dev. Allows users to subscribe to different RSS feeds, periodically fetch them, store published posts in a database, and browse through them.

More information about the assignment can be found [here](https://www.boot.dev/courses/build-blog-aggregator-golang).

## Setup
You'll need a PostgreSQL database to store information about users, the feeds they follow, and published posts.

1. Copy the provided example `.env.default` file to `.env` to ensure the correct environment variables are present when running the containers. Feel free to update the values to suit your needs.

```bash
cp .env.default .env
```

2. Spin up the PostgreSQL database container. A second container, running the image defined in the `goose` directory, will be built to execute the database migrations from `sql/schema`.

```bash
$ docker compose up
...
gator_goose  | Waiting for Postgres...
gator_goose  | Postgres is not ready yet. Sleeping 2 seconds...
...
gator_goose  | Postgres is ready. Running migrations...
gator_goose  | 2025/09/21 05:17:47 OK   001_users.sql (20.9ms)
gator_goose  | 2025/09/21 05:17:47 OK   002_feeds.sql (22.89ms)
gator_goose  | 2025/09/21 05:17:47 OK   003_feed_folllows.sql (16.38ms)
gator_goose  | 2025/09/21 05:17:47 OK   004_feeds_last_fetched_at.sql (4.22ms)
gator_goose  | 2025/09/21 05:17:47 OK   005_posts.sql (21.37ms)
gator_goose  | 2025/09/21 05:17:47 goose: successfully migrated database to version: 5
gator_goose exited with code 0
```

That's it! In the logs, you should see that the `gator_goose` container successfully completed the migrations and shut down afterwards.

```bash
$ docker ps -a
CONTAINER ID   IMAGE                 COMMAND                  CREATED          STATUS                      PORTS                                         NAMES
df47ebcac647   bootdev-gator-goose   "/entrypoint.sh"         24 seconds ago   Exited (0) 20 seconds ago                                                 gator_goose
6af06d264390   postgres:15           "docker-entrypoint.s…"   24 seconds ago   Up 23 seconds               0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp   gator_db
```

## Installation
### 1. Requires Go 1.24 or later

Use the [webi installer](https://webinstall.dev/golang/) or the [official installation instructions](https://go.dev/doc/install) to install Go. Run `go version` in your command line to make sure the installation worked.

### 2. Install the `bootdev-gator` CLI
Install the `bootdev-gator` CLI tool by running the following.

```bash
$ go install github.com/simonkosina/bootdev-gator@latest
```

Verify the installation was successful. You should see the usage guide printed in your console.

```bash
$ bootdev-gator
Usage: gator <command> [args...]

Commands:
  login <user_name>                Set the current user
  register <user_name>             Register a new user and set as current
  reset                            Reset (delete) all users
  users                            List all users
  agg <time_between_reqs>          Collect feeds every given duration (e.g. 1m, 30s)
  addfeed <feed_name> <feed_url>   Add a new feed and follow it
  feeds                            List all feeds
  follow <feed_url>                Follow a feed
  following                        List feeds you are following
  unfollow <feed_url>              Unfollow a feed
  browse [limit]                   Browse posts (default limit: 2)
```

### 3. Configuration
The project's configuration is stored in the `.gatorconfig.json` file in the user's home directory. Create the configuration file based on the example provided below, make sure to replace the `${POSTGRES_USER}` and `${POSTGRES_USER}` placeholders by the correct values from your `.env` file.

```json ~/.gatorconfig.json
{
  "db_url": "postgres://${POSTGRES_USER}:${POSTGRES_USER}@localhost:5432/gator?sslmode=disable"
}
```

## Usage

Create some users in the database.

```bash
$ bootdev-gator register simon
User was created successfully: {ID:40c16633-0148-4948-9f4e-eae992c7788c CreatedAt:2025-09-21 06:04:30.117249 +0000 +0000 UpdatedAt:2025-09-21 06:04:30.117249 +0000 +0000 Name:simon}
$ bootdev-gator register john
User was created successfully: {ID:3c4605a6-7c55-4941-9dd4-56c7caf9f087 CreatedAt:2025-09-21 06:10:14.818169 +0000 +0000 UpdatedAt:2025-09-21 06:10:14.818169 +0000 +0000 Name:john}
$ bootdev-gator users
simon
john (current)
```

Login as `simon` and add a feed.

```bash
$ bootdev-gator login simon
Current user has been set to: simon
$ bootdev-gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
Feed was added successfully: {ID:bc71133d-a0b5-46fc-9ffc-f955fd534d23 CreatedAt:2025-09-21 06:05:32.684965 +0000 +0000 UpdatedAt:2025-09-21 06:05:32.684965 +0000 +0000 Name:Boot.dev Blog Url:https://blog.boot.dev/index.xml UserID:40c16633-0148-4948-9f4e-eae992c7788c LastFetchedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}}
```

In a separate terminal, run the `agg` command to fetch added feeds every minute and save new posts.

```bash
$ bootdev-gator agg 1m
2025/09/21 08:05:44 Collecting feeds every 1m0s...
2025/09/21 08:05:44 Scraping feed 'Boot.dev Blog'
2025/09/21 08:05:44 Saving post 'The Boot.dev Beat. September 2025' (https://blog.boot.dev/news/bootdev-beat-2025-09/)
2025/09/21 08:05:44 Post 'The Boot.dev Beat. September 2025' was saved successfully
2025/09/21 08:05:46 Saving post 'Secure Random Numbers in Node.js' (https://blog.boot.dev/cryptography/node-js-random-number/)
2025/09/21 08:05:46 Post 'Secure Random Numbers in Node.js' was saved successfully
...
2025/09/21 08:05:46 Feed Boot.dev Blog collected, 381 posts found
```

Login as `john`, follow an already added feed and browse some posts. You can view existing feeds by running `feeds` command.

```bash
$ bootdev-gator login john
Current user has been set to: john
$ bootdev-gator following
$ bootdev-gator follow https://blog.boot.dev/index.xml
'john' now follow 'Boot.dev Blog' feed
$ bootdev-gator following
Boot.dev Blog
$ bootdev-gator browse
2025/09/21 08:12:10 No limit argument provided for 'browse', defaulting to 2
2025/09/21 08:12:10 Found 2 posts for user 'john'
Boot.dev Blog: The Boot.dev Beat. September 2025
  URL: https://blog.boot.dev/news/bootdev-beat-2025-09/
  Published At: Mon, 08 Sep 2025 02:00:00
  Description: <p>The training grounds are LIVE! 21,000 challenges have been generated between the launch and as I write this, and we&rsquo;re just getting started. Big things to come.</p>
  ─────────────────────────────────────────────
Boot.dev Blog: Create a Course on Boot.dev
  URL: https://blog.boot.dev/create-a-course/
  Published At: Thu, 04 Sep 2025 02:00:00
  Description: <p>We create most of our courses at Boot.dev in-house, but we also love to collaborate with talented authors! If you&rsquo;re interested in creating a course for Boot.dev, here&rsquo;s some preliminary info about how we work:</p>
```
