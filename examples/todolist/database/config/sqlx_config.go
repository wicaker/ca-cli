package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
)

// SqlxInit will connecting service to databsase using sqlx library
func SqlxInit() *sqlx.DB {
	dbHost := os.Getenv("DATABASE_HOST_MYSQL")
	dbPort := os.Getenv("DATABASE_PORT_MYSQL")
	dbUser := os.Getenv("DATABASE_USER_MYSQL")
	dbPass := os.Getenv("DATABASE_PASSWORD_MYSQL")
	dbName := os.Getenv("DATABASE_NAME_MYSQL")

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sqlx.Open(`mysql`, dsn)
	if err != nil {
		panic(err)
	}

	// defer func() {
	// 	err := dbConn.Close()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	return dbConn
}
