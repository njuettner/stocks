package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
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
	stocksFile := flag.String("file", "", "file containing stocks (one stock symbol per line)")

	flag.Parse()

	var stockSymbols []string

	if *stocksFile != "" {
		file, err := ioutil.ReadFile(*stocksFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		stockSymbols = strings.Split(string(file), "\n")
	} else {
		fmt.Println("No data file provided")
		os.Exit(1)
	}

	count := 0
	for _, stockSymbol := range stockSymbols {
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

			if count > 4 {
				// only make 5 requests per minute
				fmt.Println("Waiting 70 seconds to make new requests ...")
				time.Sleep(70 * time.Second)
				count = 0
			}

			// get the time from JSON files
			fileInfo, err := os.Stat(fmt.Sprintf("data/%s/%s.json", stockSymbol, k))
			if !os.IsNotExist(err) {
				duration := time.Since(fileInfo.ModTime())
				if duration.Hours() < 8 {
					fmt.Printf("Skipping fetching %s for %s, no new information needed ...\n", k, stockSymbol)
				}
			} else {
				count++

				response, err := http.Get(v)
				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
				} else {
					if _, err := os.Stat(fmt.Sprintf("data/%s", stockSymbol)); os.IsNotExist(err) {
						os.Mkdir(fmt.Sprintf("data/%s", stockSymbol), os.ModePerm)
					}

					data, _ := ioutil.ReadAll(response.Body)
					if k == "earnings_calendar" {
						ioutil.WriteFile(fmt.Sprintf("data/%s/%s.csv", stockSymbol, k), data, 0644)
						convertToJSON(fmt.Sprintf("data/%s/%s.csv", stockSymbol, k))
					} else {
						ioutil.WriteFile(fmt.Sprintf("data/%s/%s.json", stockSymbol, k), data, 0644)
					}
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
