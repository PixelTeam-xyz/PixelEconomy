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
             (name TEXT PRIMARY KEY, price INTEGER)`)

    Except(err, DatabaseErrorExit)

    stmt, err := db.Prepare("INSERT OR IGNORE INTO items (name, price) VALUES (?, ?)")

    Except(err, DatabaseErrorExit)

    items := []struct {
        name  string
        price int
    }{
        {"MiniVIP", 5000},
        {"VIP", 15000},
        {"MegaVIP", 35000},
        {"CustomVIP", 70000},
    }
    for _, item := range items {
        _, err = stmt.Exec(item.name, item.price)
        Except(err, DatabaseErrorExit)
    }

    _, err = db.Exec("COMMIT")
    Except(err, DatabaseErrorExit)
}
