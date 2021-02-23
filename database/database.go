package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//InitDB opens the database connection
func InitDB(dbname string, user string, pwd string) error {
	base, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", user, pwd, dbname))
	if err != nil {
		log.Panic(err)
	}
	db = base
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return nil
}
