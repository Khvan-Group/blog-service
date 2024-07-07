package db

import (
	"database/sql"
	"fmt"
	"github.com/dkhvan-dev/alabs_project/common-libraries/utils"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dbHost := utils.GetEnv("DB_HOST")
	dbName := utils.GetEnv("DB_NAME")
	dbUser := utils.GetEnv("DB_USER")
	dbPass := utils.GetEnv("DB_PASS")
	dbPort := utils.GetEnv("DB_PORT")
	sslmode := utils.GetEnv("SSLMODE")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, sslmode)

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}

	if db.Ping() != nil {
		panic(err)
	}

	DB = db
}
