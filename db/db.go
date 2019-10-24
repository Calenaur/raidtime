package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func New(user string, password string, database string) (*sql.DB, error) {
	db, err := sql.Open("mysql", user + ":" + password + "@/" + database + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	return db, nil
}