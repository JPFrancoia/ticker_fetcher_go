package yahoo

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"local/ticker_fetcher/utils"

	"github.com/tidwall/gjson"
)

const URL_API = "https://query2.finance.yahoo.com/v7/finance/options"

// Fetch the market information for one fund.
// Returns the info + pushes the info to a channel.
func FetchInfoFromYahoo(fund string, c chan<- YahooInfo, wg *sync.WaitGroup) YahooInfo {
	defer wg.Done()

	responseBody := utils.QueryApi(fmt.Sprintf("%s/%s", URL_API, fund))

	value := gjson.GetBytes(responseBody, "optionChain.result.0.quote")

	var target YahooInfo
	if err := json.Unmarshal([]byte(value.Raw), &target); err != nil {
		log.Fatal(err)
	}

	c <- target

	return target
}


type YahooInfo struct {
	Symbol        string  `json:"symbol"`
	ExchangeName  string  `json:"fullExchangeName"`
	Price         float64 `json:"regularMarketPrice"`
	PreviousClose float64 `json:"regularMarketPreviousClose"`
	Currency      string  `json:"currency"`
}

func (fi *YahooInfo) Diff() float64 {
	return (fi.Price - fi.PreviousClose) / fi.PreviousClose * 100
}
