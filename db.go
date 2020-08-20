package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

var (
	// DB
	dbURI = "postgres://uppigjzrutuobd:5c4c0f5daae8bbb4adf887f69137658840ea0ac74effdc41ed6810723ea2a9f5@ec2-54-197-254-117.compute-1.amazonaws.com:5432/ddm0qnu51ahb69"
	db, _ = sql.Open("postgres", dbURI)
)

func createTable(name string) {
	if name == "all_users" {
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + name + " (id TEXT PRIMARY KEY, balance NUMERIC);")
		if err != nil {
			log.Fatalf("[X] Could not create %s table. Reason: %s", name, err.Error())
		} else {
			log.Printf("[OK] Create %s table", name)
		}
	} else if name == "transactions" {
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS " + name + " (id TEXT, idFrom TEXT, sum NUMERIC, date INT, info TEXT);")
		if err != nil {
			log.Fatalf("[X] Could not create %s table. Reason: %s", name, err.Error())
		} else {
			log.Printf("[OK] Create %s table", name)
		}
	}
}

// Удаление таблицы
func dropTable(name string) {
	_, err := db.Exec("DROP TABLE " + name + ";")
	if err != nil {
		log.Fatalf("[X] Could not drop %s table. Reason: %s", name, err.Error())
	} else {
		log.Printf("[OK] Drop %s table", name)
	}
}

// Ищет user по id и если его нет, то создает с 0 балансом, или записывает баланс в User.Balance
func (user *User) getBalance() {
	var (
		id      string
		balance float64
	)
	rows, err := db.Query("SELECT id,balance FROM all_users;")
	defer rows.Close()
	if err != nil {
		log.Fatalf("[X] Could not select. Reason: %s", err.Error())
	} else {
		for rows.Next() {
			rows.Scan(&id, &balance)
			if id == user.ID {
				user.Balance = balance
				return
			}
		}
	}

	_, err = db.Exec("INSERT INTO all_users VALUES ('" + user.ID + "', 0);")
	if err != nil {
		log.Fatalf("[X] Could not create new user. Reason: %s", err.Error())
	} else {
		log.Printf("[OK] Add new user id=%s", user.ID)
	}
	user.Balance = 0
	return
}

func (user *User) setBalance() {
	_, err := db.Exec("UPDATE all_users SET balance = " + fmt.Sprintf("%.2f", user.Balance) + " WHERE id = '" + user.ID + "';")
	if err != nil {
		log.Fatalf("[X] Could not update balance. Reason: %s", err.Error())
	} else {
		log.Printf("[OK] Update balance for user id=%s to balance=%.2f", user.ID, user.Balance)
	}
}

func (user *User) getOperation(sort string) (Operations []Operation) {
	var (
		column string
		sc     string
	)
	switch sort {
	case "last":
		column = "date"
		sc = "ASC"
	case "new":
		column = "date"
		sc = "DESC"
	case "high":
		column = "sum"
		sc = "DESC"
	case "low":
		column = "sum"
		sc = "ASC"
	}

	var singOper Operation
	rows, err := db.Query("SELECT idFrom,sum,date,info FROM transactions WHERE id='" + user.ID + "' ORDER BY " + column + " " + sc + ";")
	defer rows.Close()
	if err != nil {
		log.Fatalf("[X] Could not select. Reason: %s", err.Error())
	} else {
		for rows.Next() {
			rows.Scan(&singOper.IDFrom, &singOper.Sum, &singOper.Date, &singOper.Info)
			Operations = append(Operations, singOper)
		}
	}
	return Operations
}

func (user *User) setOperation(idFrom string, sum string, info string) {
	_, err := db.Exec("INSERT INTO transactions VALUES ('" + user.ID + "', '" + idFrom + "','" + sum + "'," + strconv.Itoa(int(time.Now().Unix())) + ",'" + info + "');")
	if err != nil {
		log.Fatalf("[X] Could not create new operation. Reason: %s", err.Error())
	} else {
		log.Printf("[OK] Add new operation")
	}
}
