package main

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	user     = "postgres"
	password = "1111"
	host     = "localhost"
	port     = 5432
	dbname   = "postgres"
)

var doOnce sync.Once
var singleton *gorm.DB

func GetConnection() *gorm.DB {
	doOnce.Do(func() {
		connURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			user,
			password,
			host,
			port,
			dbname,
		)

		db, err := gorm.Open(
			postgres.Open(connURL),
			&gorm.Config{},
		)
		if err != nil {
			panic(err)
		}
		singleton = db
	})
	return singleton
}
