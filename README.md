# Opensearch Driver for Golang Migrate

## Introduction

This package enables migrations for managing indexes, data streams, etc., in Opensearch in a programmed and automated way.

It is built as an extension for a well-maintained and popular package for migrations in Go: [https://github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate). Since Opensearch is not an SQL database and requires all communication to be done using HTTP, this package facilitates that.

All schemas can be stored as JSON files that contain all the necessary information to make requests and communicate with Opensearch.

## Getting Started

### Installation

```shell
go get github.com/limangotech/opensearch-driver
```

## Usage

*To generate migration files, you need to install [migrate CLI](https://github.com/golang-migrate/migrate).*

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

> **You can currently use this driver only when running go-migrate from a Go file and not using the CLI.**

Example usage in a `.go` file:

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
    Addresses: []string{"https://example.com"},
    Username:  "admin",
    Password:  "admin",
  }
  
  transport, err := opensearch.NewClient(config)
  if err != nil {
    panic(err)
  }
  
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

// ...

if err = migrate.Up(); err != nil {
  panic(err)
}
```

## What is Supported

* Running migrations that respect the included [schema](pkg/schemavalidator/schema.json).
* Applying migrations both up and down.
* Validating migration files (one file or all at once) against the schema. This is especially helpful to pre-validate all migrations before applying any.
* Storing the current migration version and state in a hidden index `.migrations` in Opensearch.

## What is Not Supported

* Running go-migrate using the CLI, as it tries to open and close the driver, which is not supported due to the nature of communication with Opensearch.
* Dropping/locking/unlocking migrations due to the nature of Opensearch.
* Making requests to the Opensearch Dashboards API, which may be required if you would like to create index patterns using migrations.

## Future Plans for the Package

* Improve unit tests.
* Add integration tests.

## Have Suggestions or Improvements?

Please create an issue if any bugs are found, or if you would like to request a new feature. You can also submit a pull request with your changes if you have something to contribute.