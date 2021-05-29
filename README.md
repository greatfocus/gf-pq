# pq - A pure Go postgres driver for Go's database/sql package

## Install

	go get github.com/geratfocus/gf-pq

## Docs

<http://go.pkgdoc.org/github.com/geratfocus/gf-pq>

## Use

	package main

	import (
		"database/sql"
		"fmt"
		_ "github.com/geratfocus/pq"
		"log"
	)
	func main() {
		conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=verify-full", host, user, password, dbname)
		d, err := sql.Open("postgres", conn)
		rows, err := d.QueryRow("INSERT INTO MYTABLE(MYCOLUMN) VALUES(1),(2),(3) RETURNING ID")
		defer rows.Close()
		var ids []int
		for rows.Next() {
			var id int
			scan_err := rows.Scan(&id)
			fmt.Println(scan_err)
			if scan_err != nil {
				log.Fatal(scan_err)
			}
			ids = append(ids, id)
			fmt.Printf("id %d\n", id)
		}
	}

**Connection String Parameters**

These are a subset of the libpq connection parameters.  In addition, a
number of the [environment
variables](http://www.postgresql.org/docs/9.1/static/libpq-envars.html)
supported by libpq are also supported.  Just like libpq, these have
lower precedence than explicitly provided connection parameters.

See http://www.postgresql.org/docs/9.1/static/libpq-connect.html.

* `dbname` - The name of the database to connect to
* `user` - The user to sign in as
* `password` - The user's password
* `host` - The host to connect to. Values that start with `/` are for unix domain sockets. (default is `localhost`)
* `port` - The port to bind to. (default is `5432`)
* `sslmode` - Whether or not to use SSL (default is `require`, this is not the default for libpq)
	Valid values are:
	* `disable` - No SSL
	* `require` - Always SSL (skip verification)
	* `verify-full` - Always SSL (require verification)

See http://golang.org/pkg/database/sql to learn how to use with `pq` through the `database/sql` package.

## Tests

`go test` is used for testing.  A running PostgreSQL server is
required, with the ability to log in.  The default database to connect
to test with is "pqgotest," but it can be overridden using environment
variables.

Example:

	PGHOST=/var/run/postgresql go test pq

## Features

* SSL
* Handles bad connections for `database/sql`
* Scan `time.Time` correctly (i.e. `timestamp[tz]`, `time[tz]`, `date`)
* Scan binary blobs correctly (i.e. `bytea`)
* pq.ParseURL for converting urls to connection strings for sql.Open.
* Many libpq compatible environment variables
* Unix socket support
