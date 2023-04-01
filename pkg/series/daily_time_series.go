package series

type DailySeries struct {
	MetaData struct {
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
	} `json:"Meta Data"`
	TimeSeries map[string]struct {
		Open           string `json:"1. open"`
		High           string `json:"2. high"`
		Low            string `json:"3. low"`
		Close          string `json:"4. close"`
		Volume         string `json:"6. volume"`
		DividendAmount string `json:"7. dividend amount"`
	} `json:"Time Series (Daily)"`
}
