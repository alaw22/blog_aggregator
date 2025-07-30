# Gator

`gator` is an RSS Feed aggregator built in Go. It uses `PostgreSQL` to store 
shortened versions of posts made by websites so you can keep up with your 
favorite feeds! 


# Installation

There are a few prerequisites before we can get you up and moving.

## Go
As mentioned this application is built with Go so if you don't already have the latest version of go installed you can go to
[webi](https://webinstall.dev/golang/) for instructions on how
to download it (they have installation guides for every platform).

Once you have the go toolchain installed check to see that it
set your path environment variable appropriately by typing

```bash
go version
```

If you get an error it most likely is because you need to set your
path environment variable to include the go bin directory.

## PostgreSQL

You need to install Postgres and start up a
server. 

### Install Postgres v15 or later.

**Windows**

For windows installation go to [this website](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads) and follow the instructions

**maxOS** with brew

```zsh
brew install postgresql@15
```

**Linux / WSL (Debian)**.

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

### Make Sure the Installation Worked

```bash
psql --version
```

> You may have to add the PostgreSQL bin directory to your environment path.

### If You Are On Linux

You have to update the password:

```bash
sudo passwd postgres
```

### Start Server (macOS/Linux/Windows)

- Max: `brew services start postgresql@15`
- Linux: `sudo service postgresql start`
- Windows: service should start automatically (check services)

### Connect To Server

- Mac: `psql postgres`
- Linux: `sudo -u postgres psql`
- Windows: `psql -p <port_number> -U postgres` (if you used default port
-- 5432 -- you won't have to worry about using the `-p` option)

### Create Gator Database

Once in the psql shell you can create the gator database:

```sql
CREATE DATABASE gator;
```

To connect to the database:

```sql
\c gator
```

### Set User Password (Linux Only)

```sql
ALTER USER postgres PASSWORD 'postgres';
```

To exit the shell you can just type `exit`.

### Connection String

A connection string is an URL that has all of the information necessary to 
connect to the SQL database in question. The format is like so:

```
protocol://username:password@hostname:portnumber/database_name
```

The protocol for us would be `postgres` so an example would be:

```
postgres://postgres:postgres@localhost:5432/gator
```

Use whatever username and password you made when setting up the server. The 
default is postgres user I believe.

## Clone Repo

The best way to set up the database schema is to use the `goose` the SQL
migration application and the schema SQL files I have set up in this repository.
Make sure you have [git installed](https://git-scm.com/downloads) on your
machine if you don't already.

`cd` to a directory of your choosing to clone the repo in then:

```bash
git clone https://github.com/alaw22/gator
```

## Goose

[`goose`](https://github.com/pressly/goose) is a SQL migration application. I
have set up my `sql/schema` directory to comply with `goose` required formatting.

```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Again you can check to make sure this installation works:

```bash
goose --version
```

Once you have `goose` installed, make sure you are in the base directory
"gator" and type:

```
goose -dir sql/schema postgres <connection_string> up
```

The connection string should be in the format as aforementioned. Once you have
done this the database should have 4 tables `users`, `feeds`,
`feed_follows`, and `posts`. You can check by connecting to the database and 
then typing:

```sql
\dt
```

You should see this result:

```
gator=# \dt
              List of relations
 Schema |       Name       | Type  |  Owner
--------+------------------+-------+----------
 public | feed_follows     | table | postgres
 public | feeds            | table | postgres
 public | goose_db_version | table | postgres
 public | posts            | table | postgres
 public | users            | table | postgres
(5 rows)
```

## Install Executable

Go makes it really easy to install an application. For our purposes since we
needed the repo for the schema files we will install it locally. Make sure you
are in the base gator directory:

```bash
go install .
```

Now you can type `gator` to use the program. You will most likely get an error
because we haven't set up our `.gatorconfig.json` file just yet.

# Configuration

The only configuration we need to do is to create the `.gatorconfig.json` file.
On Windows you should create it in your `%USERPROFILE%` directory or `$HOME`
directory if on Linux or macOS.

In it you need to place your database connection URL in JSON format:

```
{
    "db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
}
```

Remember this can be whatever you set up as the username, password, and port.
The query at the end is to make sure we aren't using ssl to connect locally.


# Usage

To get started register your user!

```bash
gator register <your_username>

# sample output
# User 'bob' has been set
```

This will register you as a user in the users table and log you in.


Adding a feed, will automatically follow that feed for the logged in user:

```bash
gator addfeed <feed_name> <feed_url>

# sample output
# Feed entry made for https://techcrunch.com/feed/
# You are now following Tech Crunch
```

Aggregating feeds for logged in user:

```bash
gator agg <frequency_time_format> # ex: 10s will aggregate a feed every 10 seconds
```

The `frequency_time_format` argument must be in the format `#h#m#s` for example
if you wanted to aggregate user feeds every hour 1 minute and 9 seconds then 
you would type:

```bash
gator agg 1h1m9s
```

To browse posts from your following feeds:

```bash
gator browse [post_limit] # if no post_limit is sent then the limit is 2 posts

# sample output
# Here are the 2 most recent posts:

# Title: Our $100M Series B
# Published: 2025-07-30 13:17:28 +0000 +0000
# URL: https://oxide.computer/blog/our-100m-series-b
# Description:
#         <a href="https://news.ycombinator.com/item?id=44733817">Comments</a>

# ----------------------------------------
# Title: I launched 17 side projects. Result? I'm rich in expired domains
# Published: 2025-07-30 13:15:35 +0000 +0000
# URL: https://news.ycombinator.com/item?id=44733800
# Description:
#         <a href="https://news.ycombinator.com/item?id=44733800">Comments</a>

# ----------------------------------------
```


## Other commands

List users:

```bash
gator users

# sample output
# * bob
# * billy (current)
```

List feeds:

```bash
feeds

# sample output
# Name: Tech Crunch
# URL: https://techcrunch.com/feed/
# Username: bob
# Name: Hacker News
# URL: https://news.ycombinator.com/rss
# Username: billy
```

List feeds current user is following:

```bash
gator following

# sample output
# You are following these feeds:
#  - Hacker News
```

You can follow or unfollow a feed:

```bash
gator follow <feed_url>

# sample output
# You are now following Tech Crunch

gator unfollow <feed_url>

# sample output
# Successfully unfollowed feed: https://techcrunch.com/feed/
```



<!-- # Welcome to the Blog Aggregator
We're going to build an RSS feed aggregator in Go! We'll call it "Gator",
you know, because aggreGATOR 🐊. Anyhow, it's a CLI tool that allows users to:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the
full post


RSS feeds are a way for websites to publish updates to their content. You can
use this project to keep up with your favorite blogs, news sites, podcasts, and
more!

## Prerequisites
The project assumes that you're already familiar with the Go programming
language and SQL databases.

## Learning Goals
- Learn how to integrate a Go application with a PostgreSQL database
- Practice using your SQL skills to query and migrate a database (using sqlc
and goose, two lightweight tools for typesafe SQL in Go)
- Learn how to write a long-running service that continuously fetches new posts
from RSS feeds and stores them in the database -->