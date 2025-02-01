package main

import (
	"database/sql"
)

var stmt *sql.Stmt

func init() {
	var err error
	db, err = sql.Open("sqlite3", cnf.DatabasePath)
	Except("Failed open database (%s)", err, DatabaseErrorExit)

	_, err = db.Exec("BEGIN TRANSACTION")
	Except(err, DatabaseErrorExit)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users
             (id INTEGER PRIMARY KEY, balance INTEGER, bank INTEGER, lastWork DATETIME, rank TEXT)`)

	Except(err, DatabaseErrorExit)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items
             (name TEXT PRIMARY KEY, price INTEGER, id INTEGER)`)

	Except(err, DatabaseErrorExit)

	stmt, err := db.Prepare("INSERT OR IGNORE INTO items (name, price, id) VALUES (?, ?, ?)")

	Except(err, DatabaseErrorExit)

	items := []struct {
		name  string
		price int
		id    int64
	}{
		{"MiniVIP", 5000, 1334515871833522322},
		{"VIP", 15000, 1334515871833522323},
		{"MegaVIP", 35000, 1334515871833522324},
		{"CustomVIP", 70000, 1334515871833522325},
	}
	for _, item := range items {
		_, err = stmt.Exec(item.name, item.price, item.id)
		Except(err, DatabaseErrorExit)
	}

	_, err = db.Exec("COMMIT")
	Except(err, DatabaseErrorExit)
}
