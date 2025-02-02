package main

import (
	"database/sql"
	"fmt"
	"msg"
	"time"
)

func changeBal(user any, newBal int64) {
	stmtUpdate, err := db.Prepare("UPDATE users SET balance = ? WHERE id = ?")
	Except(err, DatabaseErrorExit)
	defer stmtUpdate.Close()

	res, err := stmtUpdate.Exec(newBal, user)
	Except(err, DatabaseErrorExit)

	rowsAffected, err := res.RowsAffected()
	Except(err, DatabaseErrorExit)

	if rowsAffected == 0 {
		stmtInsert, err := db.Prepare("INSERT INTO users (id, balance, bank, lastWork) VALUES (?, ?, ?, ?)")
		Except(err, DatabaseErrorExit)
		defer stmtInsert.Close()

		_, err = stmtInsert.Exec(user, newBal, 0, time.Now().Format("2006-01-02 15:04:05"))
		Except(err, DatabaseErrorExit)
	}
}

func changeBank(user any, newBal int64) {
	stmtUpdate, err := db.Prepare("UPDATE users SET bank = ? WHERE id = ?")
	Except(err, DatabaseErrorExit)
	defer stmtUpdate.Close()

	res, err := stmtUpdate.Exec(newBal, user)
	Except(err, DatabaseErrorExit)

	rowsAffected, err := res.RowsAffected()
	Except(err, DatabaseErrorExit)

	if rowsAffected == 0 {
		stmtInsert, err := db.Prepare("INSERT INTO users (id, bank) VALUES (?, ?)")
		Except(err, DatabaseErrorExit)
		defer stmtInsert.Close()

		_, err = stmtInsert.Exec(user, newBal)
		Except(err, DatabaseErrorExit)
	}
}

func getBal(userID any) int64 {
	var balance int64
	err := db.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0
	}
	Except(err, DatabaseErrorExit)
	return balance
}

func getBank(userID any) int64 {
	var balance int64
	err := db.QueryRow("SELECT bank FROM users WHERE id = ?", userID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0
	}
	Except(err, DatabaseErrorExit)
	return balance
}

func canWork(userID any) (bool, int) {
	var lastWork time.Time

	err := db.QueryRow("SELECT lastWork FROM users WHERE id = ?", userID).Scan(&lastWork)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, 0
		}
		msg.Fatalf("Error checking if user can work: %s", err.Error())
		return false, cnf.WorkDelay
	}

	remaining := int(float64(cnf.WorkDelay) - time.Since(lastWork).Seconds())

	if remaining <= 0 {
		return true, 0
	}

	return false, remaining
}

func getTop(x int) (topUsers []string, err error) {
	query := fmt.Sprintf("SELECT id FROM users ORDER BY (balance + bank) DESC LIMIT %d", x)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		topUsers = append(topUsers, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return topUsers, nil
}
