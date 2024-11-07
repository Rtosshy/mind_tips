package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitDBはMySQLデータベースに接続します
func InitDB() {
	var err error
	db, err = sql.Open("mysql", "root:root_password@tcp(db:3306)/blog")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

// GetDBはデータベース接続を返します
func GetDB() *sql.DB {
	return db
}
