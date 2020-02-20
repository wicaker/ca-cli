package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

// GormInit will connecting service to databsase using GORM orm
func GormInit() *gorm.DB {
	dbHost := os.Getenv("DATABASE_HOST_PG")
	// port := os.Getenv("DATABASE_PORT_PG")
	dbUser := os.Getenv("DATABASE_USER_PG")
	dbPassword := os.Getenv("DATABASE_PASSWORD_PG")
	dbName := os.Getenv("DATABASE_NAME_PG")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPassword)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic("failed to connect to database")
	}

	return db
}
