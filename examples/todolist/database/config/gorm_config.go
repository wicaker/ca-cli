package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

// GormInit will connecting service to databsase using GORM orm
func GormInit() *gorm.DB {
	dbHost := os.Getenv("DATABASE_HOST")
	// port := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPassword)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic("failed to connect to database")
	}
	fmt.Println("Connected to database using GORM orm")

	return db
}
