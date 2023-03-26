package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/njuettner/stock-data-collection/pkg/balance"
	"github.com/njuettner/stock-data-collection/pkg/cashflow"
	"github.com/njuettner/stock-data-collection/pkg/earnings"
	"github.com/njuettner/stock-data-collection/pkg/income"
	"github.com/njuettner/stock-data-collection/pkg/overview"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{symbol}/balance", GetBalance).Methods("GET")
	r.HandleFunc("/{symbol}/earnings", GetEarnings).Methods("GET")
	r.HandleFunc("/{symbol}/earnings/calendar", GetEarningsCalendar).Methods("GET")
	r.HandleFunc("/{symbol}/overview", GetOverview).Methods("GET")
	r.HandleFunc("/{symbol}/cashflow", GetCashflow).Methods("GET")
	r.HandleFunc("/{symbol}/income", GetIncome).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetBalance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params["symbol"]
	data, _ := os.ReadFile(strings.ToUpper(symbol) + "-balance.json")
	balance := balance.FinancialReport{}
	json.Unmarshal(data, &balance)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

func GetCashflow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params["symbol"]
	data, _ := os.ReadFile(strings.ToUpper(symbol) + "-cashflow.json")
	cashflow := cashflow.FinancialReport{}
	json.Unmarshal(data, &cashflow)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cashflow)
}

func GetEarnings(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params["symbol"]
	data, _ := os.ReadFile(strings.ToUpper(symbol) + "-earnings.json")
	earnings := earnings.Earnings{}
	json.Unmarshal(data, &earnings)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(earnings)
}

func GetIncome(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params["symbol"]
	data, _ := os.ReadFile(strings.ToUpper(symbol) + "-income.json")
	income := income.FinancialReport{}
	json.Unmarshal(data, &income)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(income)
}

func GetOverview(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params["symbol"]
	data, _ := os.ReadFile(strings.ToUpper(symbol) + "-overview.json")
	overview := overview.Stock{}
	json.Unmarshal(data, &overview)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(overview)
}

func GetEarningsCalendar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params["symbol"]
	data, _ := os.ReadFile(strings.ToUpper(symbol) + "-earnings_calendar.json")
	earningsCalendar := earnings.Calendar{}
	json.Unmarshal(data, &earningsCalendar)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(earningsCalendar)
}
