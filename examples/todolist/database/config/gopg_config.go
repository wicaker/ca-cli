package config

import (
	"os"

	"github.com/go-pg/pg/v9"
)

// GopgInit will connecting service to databsase using go-pg orm
func GopgInit() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DATABASE_USER"),
		Database: os.Getenv("DATABASE_NAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
	})

	return db
}
