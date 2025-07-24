# Postgres

PostgreSQL is a production-ready, open-source database. It's a great choice for
many web applications, and as a back-end engineer, it might be the single most
important database to be familiar with.

## How Does PostgreSQL Work?

Postgres, like most other database technologies, is itself a server. It listens
for requests on a port (Postgres' default is :5432), and responds to those
requests. To interact with Postgres, first you will install the server and start
it. Then, you can connect to it using a client like psql or PGAdmin.

1. Install Postgres v15 or later.

**macOS** with brew

```zsh
brew install postgresql@15
```

**Linux / WSL (Debian)**.
[Here](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql)
are the docs from Microsoft, but simply:

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

2. Ensure the installation worked. The `psql` command-line utility is the default
client for Postgres. Use it to make sure you're on version 15+ of Postgres:

```bash
psql --version
```

3. **(Linux only)** Update postgres password:

```bash
sudo passwd postgres
```

Enter a password, and be sure you won't forget it. You can just use something
easy like `postgres`.

4. Start the Postgres server in the background
 - Mac: `brew services start postgresql@15`
 - Linux: `sudo service postgresql start`

5. Connect to the server. I recommend simply using the psql client. It's the
"default" client for Postgres, and it's a great way to interact with the
database. While it's not as user-friendly as a GUI like PGAdmin, it's a great
tool to be able to do at least basic operations with.

Enter the `psql` shell:

- Mac: psql postgres
- Linux: sudo -u postgres psql

You should see a new prompt that looks like this:

```sql
postgres=#
```

6. Create a new database. I called mine gator:

```sql
CREATE DATABASE gator;
```

7. Connect to the new database:

```sql
\c gator
```

You should see a new prompt that looks like this:

```sql
gator=#
```

8. Set the user password (Linux only)

```sql
ALTER USER postgres PASSWORD 'postgres';
```

For simplicity, I used `postgres` as the password. Before, we altered the
system user's password, now we're altering the database user's password.

9. Query the database

From here you can run SQL queries against the gator database. For example, to
see the version of Postgres you're running, you can run:

```sql
SELECT version();
```

You can type `exit` to leave the `psql` shell.

# Goose Migrations
Goose is a database migration tool written in Go. It runs migrations from a set
of SQL files, making it a perfect fit for this project (we wanna stay close to
the raw SQL).

## What Is a Migration?
A migration is just a set of changes to your database table. You can have as
many migrations as needed as your requirements change over time. For example,
one migration might create a new table, one might delete a column, and one
might add 2 new columns.

An "up" migration moves the state of the database from its current schema to
the schema that you want. So, to get a "blank" database to the state it needs
to be ready to run your application, you run all the "up" migrations.

If something breaks, you can run one of the "down" migrations to revert the
database to a previous state. "Down" migrations are also used if you need to
reset a local testing database to a known state.

A "migration" in Goose is just a `.sql` file with some SQL queries and some special comments. Our first migration should just create a `users` table. The simplest format for these files is:

```
number_name.sql
```

For example, I created a file in `sql/schema` called `001_users.sql` with the
following contents:

```SQL
-- +goose Up
CREATE TABLE ...

-- +goose Down
DROP TABLE users;
```

The `-- +goose Up` and `-- +goose Down` comments are case sensitive and required.
They tell Goose how to run the migration in each direction.

## Connection String

A connection string is just a URL with all of the information needed to connect
to a database. The format is:

```
protocol://username:password@host:port/database
```

Here are examples:

- macOS (no password, your username): `postgres://wagslane:@localhost:5432/gator`
- Linux (password from last lesson, postgres user): 
`postgres://postgres:postgres@localhost:5432/gator`

It should connect you to the gator database directly. If it's working, great.
exit out of psql and save the connection string.\

# SQLC

[SQLC](https://sqlc.dev/) is an amazing Go program that generates Go code from
SQL queries. It's not exactly an [ORM](https://www.freecodecamp.org/news/what-is-an-orm-the-meaning-of-object-relational-mapping-database-tools/),
but rather a tool that makes working with raw SQL easy and type-safe.

We will use Goose to manage our database migrations (the schema), and SQLC to
generate Go code that our application can use to interact with the database
(run queries).

## Configure SQLC

You'll always run the sqlc command from the root of your project. Create a file
called sqlc.yaml in the root of your project. Here is mine:

```yaml
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
```

We're telling SQLC to look in the `sql/schema` directory for our schema structure
(which is the same set of files that Goose uses, but `sqlc` automatically ignores
"down" migrations), and in the `sql/queries` directory for queries. We're also
telling it to generate Go code in the `internal/database` directory.

Here is a query to create a user:

```sql
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;
```

`$1`, `$2`, `$3`, and `$4` are parameters that we'll be able to pass into the query in
our Go code. The `:one` at the end of the query name tells SQLC that we expect to
get back a single row (the created user).

Keep the [SQLC postgres docs](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)
handy, you'll probably need to refer to them again later.

Generate the Go code. Run `sqlc generate` from the root of your project.
It should create a new package of go code in `internal/database`.
You'll notice that the generated code relies on Google's uuid package, so
you'll need to add that to your module:

```bash
go get github.com/google/uuid
```

## Postgres driver

We need to add and import a [Postgres driver](https://github.com/lib/pq)
so our program knows how to talk to the database. Install it in your module:

```bash
go get github.com/lib/pq
```

Add this import to the top of your `main.go` file:

```go
import _ "github.com/lib/pq"
```


>This is one of my least favorite things working with SQL in Go currently.
You have to import the driver, but you don't use it directly anywhere in your
code. The underscore tells Go that you're importing it for its side effects,
not because you need to use it. - Lane\


>Be sure to migrate down and back up with goose before each run/submit, because the tests assume a clean database.