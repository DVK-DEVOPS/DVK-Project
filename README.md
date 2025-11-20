# DVK-Project - Migration branch

This branch for the DVK-Project showcases how we would implement the goose migration tool for our database. The migration in this branch targets an sqlite database since that was what was in use at the time of implementation.
This branch is to satisfy the optional subtask [**migration**](https://github.com/who-knows-inc/EK_DAT_DevOps_2025_Autumn/blob/main/00._Course_Material/01._Assignments/10._Databases_ORM_Data_scraping_Web_crawling/02._After/upgrade_the_database.md) of the [overall assignment](https://github.com/who-knows-inc/EK_DAT_DevOps_2025_Autumn/blob/main/00._Course_Material/01._Assignments/10._Databases_ORM_Data_scraping_Web_crawling/02._After/upgrade_the_database.md) of upgrading our database.

# Requirements
**goose** - https://github.com/pressly/goose 



## Usage
Review the official documentation for goose [here.](https://github.com/pressly/goose/blob/main/README.md)

Below is some cherrypicked documentation from the official documentation linked above.

# Install

```shell
go install github.com/pressly/goose/v3/cmd/goose@latest
```

This will install the `goose` binary to your `$GOPATH/bin` directory.

Binary too big? Build a lite version by excluding the drivers you don't need:

For macOS users `goose` is available as a [Homebrew
Formulae](https://formulae.brew.sh/formula/goose#default):
```shell
brew install goose
```
```
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
    validate             Check migration files without running them
```


## create

Create a new SQL migration.

    $ goose create add_some_column sql
    $ Created new file: 20170506082420_add_some_column.sql

    $ goose -s create add_some_column sql
    $ Created new file: 00001_add_some_column.sql

Edit the newly created file to define the behavior of your migration.

You can also create a Go migration, if you then invoke it with [your own goose
binary](#go-migrations):

    $ goose create fetch_user_data go
    $ Created new file: 20170506082421_fetch_user_data.go

## up

Apply all available migrations.

    $ goose up
    $ OK    001_basics.sql
    $ OK    002_next.sql
    $ OK    003_and_again.go

## up-to

Migrate up to a specific version.

    $ goose up-to 20170506082420
    $ OK    20170506082420_create_table.sql

## up-by-one

Migrate up a single migration from the current version

    $ goose up-by-one
    $ OK    20170614145246_change_type.sql

## down

Roll back a single migration from the current version.

    $ goose down
    $ OK    003_and_again.go

## down-to

Roll back migrations to a specific version.

    $ goose down-to 20170506082527
    $ OK    20170506082527_alter_column.sql

Roll back all migrations:

    $ goose down-to 0

## status
Print the status of all migrations:

    $ goose status
    $   Applied At                  Migration
    $   =======================================
    $   Sun Jan  6 11:25:03 2013 -- 001_basics.sql
    $   Sun Jan  6 11:25:03 2013 -- 002_next.sql
    $   Pending                  -- 003_and_again.go

# Migrations

goose supports migrations written in SQL or in Go.

## SQL Migrations

A sample SQL migration looks like:

```sql
-- +goose Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE post;
```

Each migration file must have exactly one `-- +goose Up` annotation. The `-- +goose Down` annotation
is optional. If the file has both annotations, then the `-- +goose Up` annotation **must** come
first.

Notice the annotations in the comments. Any statements following `-- +goose Up` will be executed as
part of a forward migration, and any statements following `-- +goose Down` will be executed as part
of a rollback.

By default, all migrations are run within a transaction. Some statements like `CREATE DATABASE`,
however, cannot be run within a transaction. You may optionally add `-- +goose NO TRANSACTION` to
the top of your migration file in order to skip transactions within that specific migration file.
Both Up and Down migrations within this file will be run without transactions.

Example usage, assuming that SQL migrations are placed in the `migrations` directory:

```go
package main

import (
    "database/sql"
    "embed"

    "github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
    var db *sql.DB
    // setup database

    goose.SetBaseFS(embedMigrations)

    if err := goose.SetDialect("postgres"); err != nil {
        panic(err)
    }

    if err := goose.Up(db, "migrations"); err != nil {
        panic(err)
    }

    // run app
}
```