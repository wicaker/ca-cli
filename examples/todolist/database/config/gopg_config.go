package config

import (
	"os"

	"github.com/go-pg/pg/v9"
)

// GopgInit will connecting service to databsase using go-pg orm
func GopgInit() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DATABASE_USER_PG"),
		Database: os.Getenv("DATABASE_NAME_PG"),
		Password: os.Getenv("DATABASE_PASSWORD_PG"),
		Addr:     os.Getenv("DATABASE_HOST_PG") + ":" + os.Getenv("DATABASE_PORT_PG"),
	})

	return db
}
