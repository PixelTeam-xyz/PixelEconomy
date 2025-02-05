package main

import (
	"database/sql"
	"errors"
	"fmt"
	"info"
	"time"
)

func changeBal(user any, newBal int64) {
	go func() {
		stmtUpdate, err := db.Prepare("UPDATE users SET balance = ? WHERE id = ?")
		Except(err, "Database error: %s")
		defer stmtUpdate.Close()

		res, err := stmtUpdate.Exec(newBal, user)
		Except(err, "Database error: %s")

		rowsAffected, err := res.RowsAffected()
		Except(err, "Database error: %s")

		if rowsAffected == 0 {
			stmtInsert, err := db.Prepare("INSERT INTO users (id, balance, bank, lastWork) VALUES (?, ?, ?, ?)")
			Except("Database error: %s", err)
			defer stmtInsert.Close()

			_, err = stmtInsert.Exec(user, newBal, 0, time.Now().Format("2006-01-02 15:04:05"))
			Except(err)
		}
	}()
}

func changeBank(user any, newBal int64) {
	go func() {
		stmtUpdate, err := db.Prepare("UPDATE users SET bank = ? WHERE id = ?")
		Except(err, "Database error: %s")
		defer stmtUpdate.Close()

		res, err := stmtUpdate.Exec(newBal, user)
		Except(err, "Database error: %s")

		rowsAffected, err := res.RowsAffected()
		Except(err, "Database error: %s")

		if rowsAffected == 0 {
			stmtInsert, err := db.Prepare("INSERT INTO users (id, bank) VALUES (?, ?)")
			Except(err, "Database error: %s")
			defer stmtInsert.Close()

			_, err = stmtInsert.Exec(user, newBal)
			Except(err, "Database error: %s")
		}
	}()
}

func getBal(userID any) int64 {
	var balance int64
	err := db.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0
	}
	Except(err, "Database error: %s")
	return balance
}

func getBank(userID any) int64 {
	var balance int64
	err := db.QueryRow("SELECT bank FROM users WHERE id = ?", userID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0
	}
	Except(err, "Database error: %s")
	return balance
}

func canWork(userID any) (bool, int) {
	var lastWork time.Time

	err := db.QueryRow("SELECT lastWork FROM users WHERE id = ?", userID).Scan(&lastWork)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, 0
		}
		info.Fatalf("Error checking if user can work: %s", err.Error())
		return false, cnf.WorkDelay
	}

	//fmt.Println("Last work:", lastWork)
	//fmt.Println("Time since last work:", time.Since(lastWork).Seconds())

	if time.Since(lastWork) > time.Duration(cnf.WorkDelay)*time.Second {
		_, updateErr := db.Exec("UPDATE users SET lastWork = ? WHERE id = ?", time.Now().Add(-time.Duration(cnf.WorkDelay)*time.Second), userID)
		if updateErr != nil {
			Except("Failed to update lastWork:", updateErr)
		} else {
			lastWork = time.Now().Add(-time.Duration(cnf.WorkDelay) * time.Second)
		}
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
