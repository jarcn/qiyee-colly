package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitMySQL() *sqlx.DB {
	var db *sqlx.DB
	var err error
	dsn := "dev:Dev_1234@tcp(172.16.5.50:3306)/jobs_data"
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	return db
}

func createTable(db *sqlx.DB, query string) error {
	_, err := db.Exec(query)
	return err
}
