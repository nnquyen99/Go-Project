package models

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"time"
)

func ConnectDataBase() *sql.DB {

	db, err := sql.Open("sqlserver", "sqlserver://sa:admin@123456@localhost:1433?database=demogo&connection+timeout=30")
	if err != nil {
		panic("Failed to connect to database")
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
