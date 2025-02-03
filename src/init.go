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

	_, err = db.Exec("COMMIT")
	Except(err, DatabaseErrorExit)
}
