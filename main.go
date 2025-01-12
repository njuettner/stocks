package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var key = os.Getenv("ALPHAVANTAGE_API_KEY")

func main() {
	stocksFile := flag.String("file", "", "file containing stocks (one stock symbol per line)")

	flag.Parse()

	var stockSymbols []string

	if *stocksFile == "" {
		fmt.Println("No data file provided")
		os.Exit(1)
	}

	file, err := os.ReadFile(*stocksFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	stockSymbols = strings.Split(string(file), "\n")

	apiLimitDaily := 0
	apiLimitPerMinute := 0
	for _, stockSymbol := range stockSymbols {
		var url = map[string]string{
			"overview":          "https://www.alphavantage.co/query?function=OVERVIEW&symbol=" + stockSymbol + "&apikey=" + key,
			"income":            "https://www.alphavantage.co/query?function=INCOME_STATEMENT&symbol=" + stockSymbol + "&apikey=" + key,
			"cashflow":          "https://www.alphavantage.co/query?function=CASH_FLOW&symbol=" + stockSymbol + "&apikey=" + key,
			"earnings":          "https://www.alphavantage.co/query?function=EARNINGS&symbol=" + stockSymbol + "&apikey=" + key,
			"balance":           "https://www.alphavantage.co/query?function=BALANCE_SHEET&symbol=" + stockSymbol + "&apikey=" + key,
			"earnings_calendar": "https://www.alphavantage.co/query?function=EARNINGS_CALENDAR&symbol=" + stockSymbol + "&horizon=12month" + "&apikey=" + key,
			"daily":             "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&outputsize=full&symbol=" + stockSymbol + "&apikey=" + key,
		}

		// Fetch stock data from Alpha Vantage API
		for k, v := range url {

			if (stockSymbol == "VGWL.DE" || stockSymbol == "VGWD.DE") && k != "daily" {
				continue
			}

			if apiLimitDaily > 24 {
				fmt.Println("Total API limit reached, aborting...")
				os.Exit(1)
			}

			if apiLimitPerMinute > 4 {
				// only make 5 requests per minute
				fmt.Println("Waiting 70 seconds to make new requests ...")
				time.Sleep(70 * time.Second)
				apiLimitPerMinute = 0
			}

			response, err := http.Get(v)
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
				continue
			}
			if _, err := os.Stat(fmt.Sprintf("data/%s", stockSymbol)); os.IsNotExist(err) {
				os.Mkdir(fmt.Sprintf("data/%s", stockSymbol), os.ModePerm)
			}

			data, _ := io.ReadAll(response.Body)
			if k == "earnings_calendar" {
				os.WriteFile(fmt.Sprintf("data/%s/%s.csv", stockSymbol, k), data, 0644)
				convertToJSON(fmt.Sprintf("data/%s/%s.csv", stockSymbol, k))
			} else {
				os.WriteFile(fmt.Sprintf("data/%s/%s.json", stockSymbol, k), data, 0644)
			}
			apiLimitPerMinute++
			apiLimitDaily++
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
