package generator_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/wicaker/cacli/fs"
	"github.com/wicaker/cacli/generator"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	expected_gopg_config = `package config

import (
	pg "github.com/go-pg/pg/v9"
	"os"
)

// GopgInit will connecting service to databsase using go-pg orm
func GopgInit() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     os.Getenv("DATABASE_HOST") + ":" + os.Getenv("DATABASE_PORT"),
		Database: os.Getenv("DATABASE_NAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		User:     os.Getenv("DATABASE_USER"),
	})

	return db
}
`
	expected_gorm_config = `package config

import (
	"fmt"
	gorm "github.com/jinzhu/gorm"
	"os"
)

// GormInit will connecting service to databsase using GORM orm
func GormInit() *gorm.DB {
	dbHost := os.Getenv("DATABASE_HOST")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPassword)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic("failed to connect to database")
	}

	return db
}
`

	expected_sql_config = `package config

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
)

// SQLInit will connecting service to databsase using database/sql standard library
func SQLInit() *sql.DB {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return dbConn
}
`

	expected_sqlx_config = `package config

import (
	"fmt"
	sqlx "github.com/jmoiron/sqlx"
	"net/url"
	"os"
)

// SqlxInit will connecting service to databsase using sqlx library
func SqlxInit() *sqlx.DB {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return dbConn
}
`
)

func TestGenerateGopgConfig(t *testing.T) {
	var (
		serviceName = "test_gopg_config"
		dirDb       = "database"
		dirConf     = "config"
		dirName     = fmt.Sprintf("%s/%s/%s", serviceName, dirDb, dirConf)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an gopg_config.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirDb)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirDb + "/" + dirConf)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate gopg_config.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGopgConfig(dirName)
		resGopg, err := newFs.FindFile(dirName + "/gopg_config.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/gopg_config.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_gopg_config, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate gopg_config file
		gen := generator.NewGeneratorService()
		err := gen.GenGopgConfig(serviceName)

		assert.Error(t, err)
	})
}

func TestGenerateGormConfig(t *testing.T) {
	var (
		serviceName = "test_gorm_config"
		dirDb       = "database"
		dirConf     = "config"
		dirName     = fmt.Sprintf("%s/%s/%s", serviceName, dirDb, dirConf)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an gorm_config.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirDb)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirDb + "/" + dirConf)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate gorm_config.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGormConfig(dirName)
		resGorm, err := newFs.FindFile(dirName + "/gorm_config.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGorm)

		data, err := ioutil.ReadFile(dirName + "/gorm_config.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_gorm_config, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate gorm_config file
		gen := generator.NewGeneratorService()
		err := gen.GenGopgConfig(serviceName)

		assert.Error(t, err)
	})
}

func TestGenerateSQLConfig(t *testing.T) {
	var (
		serviceName = "test_sql_config"
		dirDb       = "database"
		dirConf     = "config"
		dirName     = fmt.Sprintf("%s/%s/%s", serviceName, dirDb, dirConf)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an sql_config.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirDb)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirDb + "/" + dirConf)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate sql_config.go file
		gen := generator.NewGeneratorService()
		err = gen.GenSQLConfig(dirName)
		resSql, err := newFs.FindFile(dirName + "/sql_config.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resSql)

		data, err := ioutil.ReadFile(dirName + "/sql_config.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_sql_config, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate sql_config file
		gen := generator.NewGeneratorService()
		err := gen.GenGopgConfig(serviceName)

		assert.Error(t, err)
	})
}

func TestGenerateSqlxConfig(t *testing.T) {
	var (
		serviceName = "test_sqlx_config"
		dirDb       = "database"
		dirConf     = "config"
		dirName     = fmt.Sprintf("%s/%s/%s", serviceName, dirDb, dirConf)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an sqlx_config.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirDb)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirDb + "/" + dirConf)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate sqlx_config.go file
		gen := generator.NewGeneratorService()
		err = gen.GenSqlxConfig(dirName)
		resSqlx, err := newFs.FindFile(dirName + "/sqlx_config.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resSqlx)

		data, err := ioutil.ReadFile(dirName + "/sqlx_config.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_sqlx_config, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate sqlx_config file
		gen := generator.NewGeneratorService()
		err := gen.GenGopgConfig(serviceName)

		assert.Error(t, err)
	})
}
