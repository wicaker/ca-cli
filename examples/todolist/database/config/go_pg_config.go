package config

import (
	"os"
	"todolist/domain"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

// GopgInit will connecting service to databsase using go-pg orm
func GopgInit() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER_TEST"),
		Database: os.Getenv("DB_NAME_TEST"),
		Password: os.Getenv("DB_PASSWORD_TEST"),
	})

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	return db
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*domain.User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
