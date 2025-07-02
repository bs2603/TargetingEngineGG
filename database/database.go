package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/GGtargeting")
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}

}
