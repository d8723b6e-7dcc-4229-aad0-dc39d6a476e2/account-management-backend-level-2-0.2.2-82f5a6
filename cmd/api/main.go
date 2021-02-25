package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/37d7fcd0-f13b-4571-8fd3-fc12d70c7b7d/account-management-backend-level-2-0.1.5-79653d/account"
	"github.com/37d7fcd0-f13b-4571-8fd3-fc12d70c7b7d/account-management-backend-level-2-0.1.5-79653d/transaction"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := mux.NewRouter()
	// Handle healthchecks
	r.HandleFunc("/ping", healthCheckHandler)

	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Fatal(err)
	}

	transactionRepository := transaction.NewSqlLiteRepository(db)
	transactionService := transaction.NewService(transactionRepository)
	transactionHandler := transaction.NewHandler(transactionService)

	accountRepository := account.NewSqlLiteRepository(db)
	accountService := account.NewService(accountRepository, transactionService)
	accountHandler := account.NewHandler(accountService)

	r.PathPrefix("/balance").Handler(accountHandler.BalanceRouter)
	r.PathPrefix("/amount").Handler(accountHandler.AmountRouter)
	r.PathPrefix("/transaction").Handler(transactionHandler.Router)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Could not listen on port 8080")
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
