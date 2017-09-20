package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type mapRowFunc func(rows *sql.Rows)

func dbQuery(sqlQuery string, mapRow mapRowFunc) error {
	db, err := sql.Open("postgres", os.Getenv("DB_CONNECT"))
	check(err)
	defer db.Close()
	rows, err := db.Query(sqlQuery)
	check(err)
	defer rows.Close()
	for rows.Next() {
		mapRow(rows)
	}
	return rows.Err()
}
