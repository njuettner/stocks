package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var key = os.Getenv("ALPHAVANTAGE_API_KEY")

func main() {
	var stockSymbol string
	fmt.Print("Enter stock symbol: ")
	fmt.Scanln(&stockSymbol)

	var url = map[string]string{
		"overview": "https://www.alphavantage.co/query?function=OVERVIEW&symbol=" + stockSymbol + "&apikey=" + key,
		"income":   "https://www.alphavantage.co/query?function=INCOME_STATEMENT&symbol=" + stockSymbol + "&apikey=" + key,
		"cashflow": "https://www.alphavantage.co/query?function=CASH_FLOW&symbol=" + stockSymbol + "&apikey=" + key,
		"earnings": "https://www.alphavantage.co/query?function=EARNINGS&symbol=" + stockSymbol + "&apikey=" + key,
		"balance":  "https://www.alphavantage.co/query?function=BALANCE_SHEET&symbol=" + stockSymbol + "&apikey=" + key,
	}

	// Fetch stock data from Alpha Vantage API
	for k, v := range url {
		response, err := http.Get(v)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			ioutil.WriteFile(fmt.Sprintf("server/%s-%s.json", stockSymbol, k), data, 0644)
		}
	}
}
