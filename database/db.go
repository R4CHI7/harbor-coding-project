package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	HOST = "database"
	PORT = 5432
)

type db struct {
	conn *sql.DB
}

var DB db

func Init(username, password, database string) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = conn.Ping()
	if err != nil {
		panic(err)
	}

	DB = db{
		conn: conn,
	}
}
