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
             (id INTEGER PRIMARY KEY, balance REAL, bank REAL, lastWork DATETIME, lastCrime DATETIME, lastRob DATETIME)`)

	_, err = db.Exec("COMMIT")
	Except(err, DatabaseErrorExit)

	//rows, err := db.Query("SELECT id FROM users")
	//Except(err)
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var userID int
	//	err := rows.Scan(&userID)
	//	Except(err)
	//
	//	x, _ := canWork(userID)
	//	y, _ := canRob(userID)
	//	z, _ := canCrime(userID)
	//
	//	if !x || !y || !z {
	//		refresh(strconv.Itoa(userID))
	//	}
	//}
	//
	//err = rows.Err()
	//Except(err, DatabaseErrorExit)
}
