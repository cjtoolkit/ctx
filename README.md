[![](https://godoc.org/github.com/cjtoolkit/ctx?status.svg)](http://godoc.org/github.com/cjtoolkit/ctx)
[![Build Status](https://travis-ci.org/cjtoolkit/ctx.svg?branch=master)](https://travis-ci.org/cjtoolkit/ctx)

# CJToolkit Context System

Just a simple context system.

## Installation

`go get github.com/cjtoolkit/ctx`

## Usage

It's useful for storing configuration and dependencies (Dependency Injection) without having to rely on side effect.

## Example

```go
package ctx

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/cjtoolkit/ctx"
)

type Config struct {
	DbRsn string `json:"DbRsn"`
}

func GetConfig(context ctx.Context) Config {
	type ConfigContext struct{}
	return context.Persist(ConfigContext{}, func() (interface{}, error) {
		return initConfig(), nil
	}).(Config)
}

func initConfig() (config Config) {
	file, err := os.Open("setting.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func GetDatabaseConnection(context ctx.Context) *sql.DB {
	type DatabaseContext struct{}
	return context.Persist(DatabaseContext{}, func() (interface{}, error) {
		return initDatabaseConnection(context)
	}).(*sql.DB)
}

func initDatabaseConnection(context ctx.Context) (*sql.DB, error) {
	return sql.Open("postgres", GetConfig(context).DbRsn)
}
```
