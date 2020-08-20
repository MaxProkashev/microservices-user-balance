package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	addr = flag.String("addr", ":9000", "http service address")
)

func main() {
	// предварительная настройка порт, бд
	flag.Parse()
	dropTable("all_users")
	dropTable("transactions")
	createTable("all_users")
	createTable("transactions")

	r := mux.NewRouter()

	r.HandleFunc("/user/{id:[0-9]+}/{act:(?:add|del)}", userHandle).Methods("GET").Queries("sum", "{sum:[0-9]+}")

	r.HandleFunc("/user/transfer", userTransfer).Methods("GET").Queries("sum", "{sum:[0-9]+}", "idFrom", "{idFrom:[0-9]+}", "idTo", "{idTo:[0-9]+}")

	r.HandleFunc("/user/{id:[0-9]+}/balance", userBalance).Methods("GET").Queries("currency", "{currency}")
	r.HandleFunc("/user/{id:[0-9]+}/balance", userBalanceRUB).Methods("GET")

	r.HandleFunc("/user/{id:[0-9]+}/transactions", userTransactions).Methods("GET").Queries("sort", "{sort:last|new|high|low}")

	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
