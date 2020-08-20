package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// User - пользователь с
// ID (string)
// Balance (int)
type User struct {
	ID       string  `json:"id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

// Operation - информация об одной операции
type Operation struct {
	IDFrom string  `json:"idFrom"`
	Sum    float64 `json:"sum"`
	Date   int     `json:"date"`
	Info   string  `json:"info"`
}
type singLine struct {
	Num           int       `json:"num"`
	SingOperation Operation `json:"operation"`
}
type respSort struct {
	Line []singLine `json:"line"`
}

func userHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sum, _ := strconv.ParseFloat(vars["sum"], 64)

	if sum == 0 {
		http.Error(w, "422 The amount is outside the acceptable range for the store", 422)
		return
	}

	user := User{
		ID:       vars["id"],
		Currency: "RUB",
	}
	user.getBalance() // баланс который, был записывается в user.balance в RUB

	switch vars["act"] {
	case "add":
		user.Balance += sum
		user.setBalance() // новый баланс
		user.setOperation(user.ID, vars["sum"], "user has deposited money")
	case "del":
		if user.Balance < sum {
			http.Error(w, "608 Insufficient funds for card transaction", 608)
			return
		}
		user.Balance -= sum
		user.setBalance() // новый баланс
		user.setOperation(user.ID, vars["sum"], "user has withdrawn money from the account")
	default:
		http.Error(w, "400 Bad Syntax act=add|del", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func userTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sum, _ := strconv.ParseFloat(vars["sum"], 64)

	if sum == 0 {
		http.Error(w, "422 The amount is outside the acceptable range for the store", 422)
		return
	}

	userFrom := User{
		ID:       vars["idFrom"],
		Currency: "RUB",
	}
	userTo := User{
		ID:       vars["idTo"],
		Currency: "RUB",
	}
	userFrom.getBalance() // баланс который был, записывается в userFrom.balance в RUB
	userTo.getBalance()   // баланс который был, записывается в userTo.balance в RUB
	if userFrom.Balance < sum {
		http.Error(w, "608 Insufficient funds for card transaction", 608)
		return
	}

	userFrom.Balance -= sum
	userTo.Balance += sum
	userFrom.setBalance()
	userTo.setBalance()

	userFrom.setOperation(userTo.ID, vars["sum"], "user transferred money")
	userTo.setOperation(userFrom.ID, vars["sum"], "user transferred money")

	json.NewEncoder(w).Encode(userFrom)
	json.NewEncoder(w).Encode(userTo)
}

func userBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	user := User{
		ID:       vars["id"],
		Currency: vars["currency"],
	}

	user.getBalance() // баланс который был, записывается в user.balance в RUB

	resp, err := http.Get("https://api.exchangeratesapi.io/latest?base=RUB&symbols=" + user.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	var info map[string]map[string]float64
	json.NewDecoder(resp.Body).Decode(&info)

	rate := info["rates"][user.Currency]
	if rate != 0 {
		user.Balance *= rate
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "400 No data available for this currency", http.StatusBadRequest)
	}
}

func userBalanceRUB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := User{
		ID:       mux.Vars(r)["id"],
		Currency: "RUB",
	}
	user.getBalance()

	json.NewEncoder(w).Encode(user)
}

func userTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := User{
		ID:       mux.Vars(r)["id"],
		Currency: "RUB",
	}
	user.getBalance()

	var resp []singLine
	var temp singLine

	for temp.Num, temp.SingOperation = range user.getOperation(mux.Vars(r)["sort"]) {
		resp = append(resp, temp)
	}

	json.NewEncoder(w).Encode(resp)
}
