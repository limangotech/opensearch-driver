# Opensearch driver for golang migrate

## Introduction

My team and I had a task to introduce migration for managing indexes, data stream, etc in Opensearch, and we quickly
realized that there is no ready solutions for that. One option was to manage manually but was not an option for us
because all other parts were automated and repeatable.

After some discussion we decided to use well-maintained and popular package for migrations in
go: [https://github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate). The only
problem we faced back then was that it had out-of-the-box support for SQL drivers, but not Opensearch which is not an
SQL database, and it would require all communication to be done using HTTP communication.

Thanks to migrate package's elasticity we decided to implement a driver to communicate with Opensearch and came up with
our own schema for migrations: we decided to go with JSON files that contain all necessary infos to be able to make
requests and communicate with Opensearch.

## Getting started

### Installation

``` shell 
go get github.com/limangotech/opensearch-driver
```

## Usage

*To generate migration files you would need to install [migrate cli](https://github.com/golang-migrate/migrate).*

Example migration creation:

```shell
migrate create -ext json -dir db/migrations -seq create_data_stream_template
```

---

Example migration **up** file `db/migrations/000001_create_data_stream_template.up.json`:

```json
{
  "method": "PUT",
  "url": "/_index_template/logs",
  "params": {
    "create": [
      "true"
    ]
  },
  "body": {
    "index_patterns": [
      "logs*"
    ],
    "data_stream": {
      "timestamp_field": {
        "name": "@timestamp"
      }
    },
    "template": {
      "mappings": {
        "properties": {
          "arbitrary": {
            "type": "flat_object"
          }
        }
      }
    },
    "priority": 100
  }
}
```

Example migration **down** file `db/migrations/000001_create_data_stream_template.down.json`:

```json
{
  "method": "DELETE",
  "url": "/_index_template/logs"
}
```

---

> **You can currently use this driver only when running go-migrate from go file and not using cli.**
>
Exmaple usage in *.go file:

```go
package main

import (
	opensearchdriver "github.com/limangotech/opensearch-driver/pkg/opensearch"
	"github.com/opensearch-project/opensearch-go/v2"
)

const (
	sourceDir = "db/migrations"
	dbName    = "opensearch"
)

func main() {
	config := opensearch.Config{
		Addresses: "https://example.com",
		Username:  "admin",
		Password:  "admin",
	}
	transport, err := opensearch.NewClient(config)
	manager := opensearchdriver.NewMigrationsIndexManager(transport)
	driver := opensearchdriver.NewDriver(transport, manager)

	migrate, err := gomigrate.NewWithDatabaseInstance("file://"+sourceDir, dbName, driver)
	if err != nil {
		panic(err)
	}

	if err = migrate.Up(); err != nil {
		panic(err)
	}
}
```

---

Example of pre-validation of all migration files before actually applying any:

```go
validator := opensearchdriver.NewSchemaValidator()
if err := validator.Dir(sourceDir); err != nil {
panic(err)
}

...

if err = migrate.Up(); err != nil {
panic(err)
}

```

## What is supported:

* Running migrations that respect included [schema](pkg/schemavalidator/schema.json).
* Applying migrations both up and down.
* Validating migrations (one file or all at once) files against schema. This is especially helpful to pre-validate all
  migrations before applying any.
* Storing current migration version and state in a hidden index `.migrations` in Opensearch.

## What is not supported:

* Running go-migrate using cli as it is trying to open and close driver which is not supported due to the nature of
  communication with Opensearch.
* Dropping/Locking/Unlocking again due to the nature of Opensearch.
* Making requests to Opensearch Dasboards API, this can be required if you would like to create index patterns using
  migrations.

## Future plans for the package

* Improve unit tests
* Add integration tests

## Have suggestions or improvements?

I hope that this package will become better with the help of the community. So please create a pull request with your
changes if you have something to contribute.

It is worth mentioning that at this stage it fulfills all the requirements my team and I had to implement migrations in
Opensearch.