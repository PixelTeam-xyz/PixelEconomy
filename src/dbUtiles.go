package main

import (
    "database/sql"
    "errors"
    "fmt"
    dsc "github.com/bwmarrin/discordgo"
    "info"
    "time"
)

func changeBal(user any, newBal float64) {
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

func changeBank(user any, newBal float64) {
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

func getBal(userID any) (balance float64) {
    err := db.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&balance)
    if err == sql.ErrNoRows {
        return 0
    }
    Except(err, "Database error: %s")
    return balance
}

func getBank(userID any) (balance float64) {
    err := db.QueryRow("SELECT bank FROM users WHERE id = ?", userID).Scan(&balance)
    if err == sql.ErrNoRows {
        return 0
    }
    Except(err, "Database error: %s")
    return balance
}

func canWork(userID any) (bool, int) {
    return canPerformAction(userID, "lastWork", cnf.WorkDelay)
}

func canCrime(userID any) (bool, int) {
    return canPerformAction(userID, "lastCrime", cnf.CrimeDelay)
}

func canRob(userID any) (bool, int) {
    return canPerformAction(userID, "lastRob", cnf.RobDelay)
}

func canPerformAction(userID any, action string, delay int) (bool, int) {
    var lastActionTime time.Time
    query := fmt.Sprintf("SELECT %s FROM users WHERE id = ?", action)

    err := db.QueryRow(query, userID).Scan(&lastActionTime)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return true, 0
        }
        info.Fatalf("Error checking if user can %s: %s", action, err.Error())
        return false, delay
    }

    if time.Since(lastActionTime) > time.Duration(delay)*time.Second {
        updateQuery := fmt.Sprintf("UPDATE users SET %s = ? WHERE id = ?", action)
        _, updateErr := db.Exec(updateQuery, time.Now().Add(-time.Duration(delay)*time.Second), userID)
        if updateErr != nil {
            Except(fmt.Sprintf("Failed to update %s: ", action)+"%s", updateErr)
        }
    }

    remaining := int(float64(delay) - time.Since(lastActionTime).Seconds())

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

func getMultiplier(userID string, msg *dsc.MessageCreate) (multi float64) {
    multi = 1
    for _, item := range items {
        if HasRole(msg.GuildID, userID, item.RoleID) && item.Multiplier > multi {
            multi = item.Multiplier
        }
    }
    return
}
