package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var key = os.Getenv("ALPHAVANTAGE_API_KEY")

func main() {
	var stockSymbol string
	fmt.Print("Enter stock symbol: ")
	fmt.Scanln(&stockSymbol)

	var url = map[string]string{
		"overview":          "https://www.alphavantage.co/query?function=OVERVIEW&symbol=" + stockSymbol + "&apikey=" + key,
		"income":            "https://www.alphavantage.co/query?function=INCOME_STATEMENT&symbol=" + stockSymbol + "&apikey=" + key,
		"cashflow":          "https://www.alphavantage.co/query?function=CASH_FLOW&symbol=" + stockSymbol + "&apikey=" + key,
		"earnings":          "https://www.alphavantage.co/query?function=EARNINGS&symbol=" + stockSymbol + "&apikey=" + key,
		"balance":           "https://www.alphavantage.co/query?function=BALANCE_SHEET&symbol=" + stockSymbol + "&apikey=" + key,
		"earnings_calendar": "https://www.alphavantage.co/query?function=EARNINGS_CALENDAR&symbol=" + stockSymbol + "&horizon=12month" + "&apikey=" + key,
	}

	// Fetch stock data from Alpha Vantage API
	for k, v := range url {

		// get the time from JSON files
		fileInfo, err := os.Stat(fmt.Sprintf("server/%s-%s.json", stockSymbol, k))
		if !os.IsNotExist(err) {
			duration := time.Since(fileInfo.ModTime())
			// Check if the duration is greater than a year
			if duration.Hours() < 24*365 && k == "overview" {
				fmt.Printf("Skipping fetching %s for %s, no new information needed ...\n", k, stockSymbol)
				// Check if the duration is greater than a month
			} else if duration.Hours() < 24*30 {
				fmt.Printf("Skipping fetching %s for %s, no new information needed ...\n", k, stockSymbol)
			}
		} else {
			response, err := http.Get(v)
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				if k == "earnings_calendar" {
					ioutil.WriteFile(fmt.Sprintf("server/%s-%s.csv", stockSymbol, k), data, 0644)
					convertToJSON(fmt.Sprintf("server/%s-%s.csv", stockSymbol, k))
				} else {
					ioutil.WriteFile(fmt.Sprintf("server/%s-%s.json", stockSymbol, k), data, 0644)
				}
			}
		}
	}
}

func convertToJSON(fileName string) {
	csvFile, _ := os.Open(fileName)

	reader := csv.NewReader(csvFile)
	header, _ := reader.Read()

	rows := make([]map[string]string, 0)

	// Read the CSV rows and convert them to maps
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		// Create a new map to hold the row data
		data := make(map[string]string)

		// Add each column to the map
		for i := 0; i < len(header); i++ {
			data[header[i]] = row[i]
		}

		// Add the map to the slice of rows
		rows = append(rows, data)
	}

	// Convert the slice of maps to a JSON string
	jsonString, _ := json.Marshal(rows)

	jsonFile, _ := os.Create(fmt.Sprintf("%s", strings.Replace(fileName, ".csv", ".json", 1)))
	defer jsonFile.Close()

	_, _ = jsonFile.Write(jsonString)
}
