package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//InitDB opens the database connection
func InitDB(dbname string, user string, pwd string) error {
	s := fmt.Sprintf("%s:%s@/%s", user, pwd, dbname)
	base, err := sql.Open("mysql", s)
	if err != nil {
		panic(err)
	}
	db = base
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	groupTokens = make(map[string]int)
	tokensToUsername = make(map[string]string)
	return nil
}

//CloseDB closes the connection to the database
func CloseDB() {
	db.Close()
}
